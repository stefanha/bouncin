"""Main bouncer code"""

import os
import ConfigParser
import logging
import reloader
import events
import runloop

__all__ = []

def get_random_port():
    return 10000 # TODO

def get_default_config():
    section = 'bouncer'
    defaults = {
        'modules': '', # TODO
        'port_range_start': get_random_port(),
    }

    cfg = ConfigParser.RawConfigParser()
    cfg.add_section(section)
    for option, value in defaults.iteritems():
        cfg.set(section, option, value)
    return cfg

def commit_config(cfg):
    config_path = os.path.expanduser('~/.bouncin')
    config_tmp_path = config_path + '.tmp'

    # The config file should be private
    old_umask = os.umask(0077)

    # Attempt to update the config file safely in the face of crashes
    try:
        f = open(config_tmp_path, 'w')
        cfg.write(f)
        f.flush()
        if hasattr(os, 'fdatasync'):
            os.fdatasync(f.fileno())
        f.close()
        try:
            os.rename(config_tmp_path, config_path)
        except OSError:
            # On Windows rename fails if file exists, on POSIX the rename
            # replaces the file and is atomic.
            os.remove(config_path)
            os.rename(config_tmp_path, config_path)
    except (OSError, IOError), e:
        logging.error('commit_config: ' + str(e))

    os.umask(old_umask)

cfg = get_default_config()
commit_config(cfg)

try:
    runloop.run()
except reloader.Reload:
    # Notify so modules can stash state before being reloaded
    events.notify('pre-reload')
    raise
