### ECM Server
This directory contains source code for ECM server. You can start it as
following:
```

# creeate server config file, e.g. server_config.json
# adjust your vault area to where your actual vault is
{
    "port": 8888,
    "vault_area": "/path/.ecm",
    "verbose": 1
}

# start the server
./ecm -server server_config.json
```

The ECM server support the following list of APIs
- GET URL/Vault provides list of records
- GET URL/Vault/recordID provides encrypted data record
- DELETE URL/Vault/recordID delete data record
- POST URL/Vault -d payload, upload record to the server
- GET URL/Vault/token provides token to use in API requests

For example,
```
# to get records from the vault named Primary
curl http;//localhost:8888/vault/Primary

# to get specific record content from vault Primary:
curl http;//localhost:8888/vault/Primary/fb26fd73-ea17-49f5-b38b-cf17575f1264.aes

# to post record to the vault Primary:
curl -X POST -d@your_record.json http;//localhost:8888/vault/Primary

# to delete record from the vault Primary
curl -X DELETE http;//localhost:8888/vault/Primary/fb26fd73-ea17-49f5-b38b-cf17575f1264.aes 

```
