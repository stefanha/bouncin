"""Configuration system"""

import os
import ConfigParser

__all__ = ['commit']

CONFIG_PATH = os.path.expanduser('~/.bouncin')

def _get_random_port():
    return 10000 # TODO

def _get_default_config():
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

_cfg = _get_default_config()

# TODO load config from file

def commit():
    """Rewrite configuration file on disk"""

    config_tmp_path = CONFIG_PATH + '.tmp'

    # The config file should be private
    old_umask = os.umask(0077)

    # Attempt to update the config file safely in the face of crashes
    try:
        f = open(config_tmp_path, 'w')
        _cfg.write(f)
        f.flush()
        if hasattr(os, 'fdatasync'):
            os.fdatasync(f.fileno())
        f.close()
        try:
            os.rename(config_tmp_path, CONFIG_PATH)
        except OSError:
            # On Windows rename fails if file exists, on POSIX the rename
            # replaces the file and is atomic.
            os.remove(CONFIG_PATH)
            os.rename(config_tmp_path, CONFIG_PATH)
    except (OSError, IOError), e:
        logging.error('%s: %s' % (__name__, str(e)))

    os.umask(old_umask)
