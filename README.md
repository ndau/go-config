# go-config
Provides a standard configuration interface for Ndau Go projects

# Usage
Put your configuration file under the folder specified by the environment variables `NDAU_CONFIG_NAME` and `NDAU_CONFIG_PATH`

# Example
```sh
export NDAU_CONFIG_NAME=config
export NDAU_CONFIG_PATH=./config 
```

This is a sample of `config/yaml` file:
```
env:
    DB_CONNECTION_STRING: 'postgress://user_name:password@host:5432/postgres'
    SOME_SETTING: value
```

