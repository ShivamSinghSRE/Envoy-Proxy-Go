#!/bin/bash

# Variables
SOURCE_VAULT="sourcekeyvault"
TARGET_VAULT="targetkeyvault"

# Get the list of secrets from the source Key Vault
secrets=$(az keyvault secret list --vault-name $SOURCE_VAULT --query "[].id" -o tsv)

# Loop through each secret and copy it to the target Key Vault
for secret_id in $secrets; do
    secret_name=$(basename $secret_id)
    secret_value=$(az keyvault secret show --vault-name $SOURCE_VAULT --name $secret_name --query "value" -o tsv)

    # Set the secret in the target Key Vault
    az keyvault secret set --vault-name $TARGET_VAULT --name $secret_name --value "$secret_value"

    echo "Secret '$secret_name' copied from '$SOURCE_VAULT' to '$TARGET_VAULT'"
done

echo "All secrets have been copied successfully."
~
