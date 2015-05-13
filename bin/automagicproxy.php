<?php

/*
 * Automate haproxy to use multiple webservers based on container name
 * ==============================================================================
 * NOTES:
 * - Currently, should only be used in dev.
 * - No real error checking. Will come in later versions depending on it's usage.
 */

/* Setup the timezone for logging */
# TODO: Update base container to have the correct timezone
$timezone = getenv('AUTOMAGIC_TIMEZONE') ? getenv('AUTOMAGIC_TIMEZONE') : 'America/Denver';
date_default_timezone_set($timezone);

$timeout = getenv('AUTOMAGIC_TIMEOUT') ? getenv('AUTOMAGIC_TIMEOUT') : 10000;

_log('Loading templates');
$t_header = file_get_contents(__DIR__ . '/templates/header.cfg');
$t_header = str_replace('__timeout__', $timeout, $t_header);
_log('Timeout was set to: ' . $timeout);
$t_frontend = file_get_contents(__DIR__ . '/templates/frontend.cfg');
$t_backend = file_get_contents(__DIR__ . '/templates/backend.cfg');
_log('Templates loaded');

/* Grab additional ports to use besides port 80 */
$ports = getenv('AUTOMAGIC_PORTS') ? explode(',', getenv('AUTOMAGIC_PORTS')) : array();
$ports[] = 80;
$ports = array_map('intval', $ports);
foreach($ports as $p){
    _log('Will use ' . $p . ' as a valid http private port');
}

/* Save a spot for our configuration */
$web_containers = array();
$frontend = $backend = '';

/* Giddy Up */
if($containers = getContainers()){
    _log('Found ' . count($containers) . ' running containers');
    /* Go through each container and see if it qualifies as a web container */
    foreach($containers as $c){
        $c_name = ltrim($c['Names'][0],'/');
        foreach($c['Ports'] as $c_port){
            if(isset($c_port['PrivatePort']) && in_array($c_port['PrivatePort'], $ports)){
                if(isset($c_port['PublicPort'])){
                    $web_containers[$c_name] = $c_port['PublicPort'];
                    _log($c_name . '* -> 172.17.42.1:' . $c_port['PublicPort'], 'WEB');
                }
            }
        }
    }

    /* Now we have the containers, lets see if we were given any configuration files */
    $mappings_dir = __DIR__ . '/../config/';
    if(file_exists($mappings_dir)) {
        $mappings = scandir($mappings_dir);
        foreach($mappings as $mapping) {
            if($mapping == '.' || $mapping == '..'){
                continue;
            } else {
                $contents = file_get_contents($mappings_dir . $mapping);
                $contents = explode("\n", $contents);
                foreach($contents as $map) {
                    if(strpos($map, ':') === FALSE) {
                        continue;
                    } else {
                        list($incoming, $container) = explode(':', $map);
                        if(isset($web_containers[$container])){
                            $web_containers[$incoming] = $web_containers[$container];
                            _log($incoming . '* -> 172.17.42.1:' . $web_containers[$container], 'WEB');
                        }
                    }
                }
            }
        }
    }

    /**
     * Order the containers descending by name length.
     * We want shorter names to appear last so if multiple containers match the shortest one is used.
     */
    $keys = array_map('strlen', array_keys($web_containers));
    array_multisort($keys, SORT_DESC, $web_containers);

    /* Ok, now that we have our containers, lets start to write our configuration*/
    foreach($web_containers as $wc_host=>$wc_port){
        $frontend .= str_replace('__host__', $wc_host, $t_frontend);
        $be_temp  = str_replace('__port__', $wc_port, $t_backend);
        $be_temp  = str_replace('__host__', $wc_host, $be_temp);
        $backend .= $be_temp;
    }

    $new_config = $t_header;
    $new_config = str_replace('__frontend__', $frontend, $new_config);
    $new_config = str_replace('__backend__', $backend, $new_config);

    if($config_file = writeConfig($new_config)){
        restartHAPROXY($config_file);
    }

} else {
    _log('Unable to fetch containers', 'ERROR', TRUE);
}

/* Quick functions */
function getContainers(){
    $api = 'http://172.17.42.1:2375/containers/json';
    /* Debugging with boot2docker */
    //$api = 'http://192.168.59.103:2375/containers/json';
    _log('Fetching containers using Docker remote API: ' . $api );
    _log('Is socat running? If not try: $(docker run sequenceiq/socat)', 'QUESTION');
    $containers = json_decode(trim(`curl -s $api`), TRUE);
    return is_array($containers) ? $containers : _log('A problem occured with the Docker Remote API', 'ERROR', TRUE);
}

function _log($string, $type = 'INFO', $die = false){
    echo date("F j, Y, g:i a") . "   " . '['. $type . ']' . "\t\t" . $string . PHP_EOL;
}


function restartHAPROXY($new_config){
    $results = `haproxy -f $new_config -D -p /var/run/haproxy.pid -sf $(cat /var/run/haproxy.pid)`;
    _log("haproxy -f $new_config -D -p /var/run/haproxy.pid -sf $(cat /var/run/haproxy.pid)");
}

function writeConfig($new_config){
    $file_name = __DIR__ . '/config/haproxy.cfg';
    $results = file_put_contents($file_name, $new_config);
    if($results !== FALSE){
        _log($file_name . ' written to succesfully!', 'SUCCESS');
    }
    return $results === FALSE ? FALSE : $file_name;
}
