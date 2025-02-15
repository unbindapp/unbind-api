#!/bin/bash
set -e

# Function to base64url encode
base64url() {
    base64 | tr '+/' '-_' | tr -d '='
}

# Wait for Zitadel to be ready
wait_for_zitadel() {
    echo "Waiting for Zitadel to be ready..."
    until curl -s http://zitadel:8080/healthz | grep -q "SERVING"; do
      echo "Waiting for Zitadel..."
      sleep 2
    done
    echo "Zitadel is ready!"
}

# Function to get access token using service account
get_access_token() {
    local SA_KEY_FILE="./.data/zitadel/zitadel-admin-sa.json"
    if [ ! -f "$SA_KEY_FILE" ]; then
        echo "Service account key file not found at $SA_KEY_FILE"
        exit 1
    fi

    echo "Reading service account key file..."
    cat "$SA_KEY_FILE"

    # Extract necessary fields from the JSON key file
    local KEY_ID=$(jq -r .keyId "$SA_KEY_FILE")
    if [ -z "$KEY_ID" ]; then
        echo "Failed to extract keyId from service account file"
        exit 1
    fi
    echo "Found key ID: $KEY_ID"

    local KEY=$(jq -r .key "$SA_KEY_FILE")
    if [ -z "$KEY" ]; then
        echo "Failed to extract key from service account file"
        exit 1
    fi
    echo "Found private key"

    # Create JWT token
    local NOW=$(date +%s)
    local EXP=$((NOW + 3600))
    
    local HEADER='{"alg":"RS256","kid":"'$KEY_ID'"}'
    local PAYLOAD='{"aud":["zitadel"],"exp":'$EXP',"iat":'$NOW',"iss":"zitadel-admin-sa","sub":"zitadel-admin-sa","scope":["openid","profile","email"]}'
    
    echo "Created header: $HEADER"
    echo "Created payload: $PAYLOAD"

    # Base64URL encode header and payload
    local B64_HEADER=$(echo -n "$HEADER" | base64url)
    local B64_PAYLOAD=$(echo -n "$PAYLOAD" | base64url)
    
    # Create signature
    echo "Signing JWT..."
    echo "$KEY" > /tmp/private.key
    local SIGNATURE=$(echo -n "$B64_HEADER.$B64_PAYLOAD" | openssl dgst -sha256 -sign /tmp/private.key | base64url)
    rm /tmp/private.key
    
    local JWT="$B64_HEADER.$B64_PAYLOAD.$SIGNATURE"
    echo "Created JWT: $JWT"

    # Get the access token
    echo "Requesting access token from Zitadel..."
    local RESPONSE=$(curl -v -X POST "http://zitadel:8080/oauth/v2/token" \
        -H "Content-Type: application/x-www-form-urlencoded" \
        -d "grant_type=urn:ietf:params:oauth:grant-type:jwt-bearer" \
        -d "assertion=$JWT")
    
    echo "Token response: $RESPONSE"
    local ACCESS_TOKEN=$(echo "$RESPONSE" | jq -r .access_token)
    
    if [ -z "$ACCESS_TOKEN" ] || [ "$ACCESS_TOKEN" = "null" ]; then
        echo "Failed to get access token"
        exit 1
    fi
    
    echo "$ACCESS_TOKEN"
}

# Create project and OIDC configuration
setup_project() {
    local ACCESS_TOKEN=$1
    
    # Create project
    echo "Creating Unbind API project..."
    local PROJECT_RESPONSE=$(curl -s -X POST "http://zitadel:8080/oauth/v2/projects" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -H "Content-Type: application/json" \
        -d '{
            "name": "Unbind API",
            "projectRoleAssertion": true,
            "projectRoleCheck": true
        }')
    
    echo "Project creation response: $PROJECT_RESPONSE"
    local PROJECT_ID=$(echo "$PROJECT_RESPONSE" | jq -r .id)
    
    if [ -z "$PROJECT_ID" ] || [ "$PROJECT_ID" = "null" ]; then
        echo "Failed to create project"
        exit 1
    fi
    
    # Create OIDC application
    echo "Creating OIDC application..."
    local APP_RESPONSE=$(curl -s -X POST "http://zitadel:8080/oauth/v2/projects/$PROJECT_ID/apps/oidc" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -H "Content-Type: application/json" \
        -d '{
            "name": "Unbind Local Dev",
            "redirectUris": ["http://localhost:8089/auth/callback"],
            "responseTypes": ["CODE"],
            "grantTypes": ["AUTHORIZATION_CODE", "REFRESH_TOKEN"],
            "appType": "USER_AGENT",
            "authMethodType": "NONE",
            "version": "V2"
        }')
    
    echo "App creation response: $APP_RESPONSE"
    
    # Save OIDC configuration
    local CLIENT_ID=$(echo "$APP_RESPONSE" | jq -r .clientId)
    local CLIENT_SECRET=$(echo "$APP_RESPONSE" | jq -r .clientSecret)
    
    if [ -z "$CLIENT_ID" ] || [ "$CLIENT_ID" = "null" ]; then
        echo "Failed to create OIDC application"
        exit 1
    fi
    
    echo "Saving OIDC configuration..."
    echo '{
        "issuer": "http://localhost:8082",
        "clientId": "'$CLIENT_ID'",
        "clientSecret": "'$CLIENT_SECRET'",
        "redirectUri": "http://localhost:8089/auth/callback"
    }' > ./.data/zitadel/oidc.json
    
    echo "OIDC configuration saved to ./.data/zitadel/oidc.json"
}

main() {
    wait_for_zitadel
    
    echo "Getting access token..."
    ACCESS_TOKEN=$(get_access_token)
    
    if [ -z "$ACCESS_TOKEN" ]; then
        echo "Failed to get access token"
        exit 1
    fi
    
    setup_project "$ACCESS_TOKEN"
    
    echo "Zitadel initialization completed successfully!"
}

main