#!/bin/bash

# SOPS Helper Script for Naytife Platform
# This script helps manage encrypted secrets across environments

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"
SECRETS_DIR="${PROJECT_ROOT}/deploy/secrets"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

usage() {
    echo -e "${BLUE}SOPS Helper Script for Naytife Platform${NC}"
    echo ""
    echo "Usage: $0 [COMMAND] [ENVIRONMENT] [SECRET_FILE]"
    echo ""
    echo "Commands:"
    echo "  decrypt    Decrypt a secret file for editing"
    echo "  encrypt    Encrypt a secret file after editing"
    echo "  edit       Edit a secret file (decrypt, edit, encrypt)"
    echo "  view       View decrypted content of a secret file"
    echo "  validate   Validate all encrypted secret files"
    echo "  list       List all secret files for an environment"
    echo "  doctor     Diagnose key/recipient issues for a secret file"
    echo "  updatekeys Re-wrap one or all secrets with current recipients"
    echo "  keygen     Generate new age key"
    echo "  keyshow    Show public key" 
    echo "  keycheck   Check key security and permissions"
    echo ""
    echo "Environments:"
    echo "  local      Local development environment"
    echo "  staging    Staging environment"
    echo "  production Production environment"
    echo ""
    echo "Examples:"
    echo "  $0 edit local postgres-secret.yaml"
    echo "  $0 view staging backend-secret.yaml"
    echo "  $0 validate"
    echo "  $0 list production"
    echo "  $0 keygen"
    echo "  $0 keyshow"
    echo ""
}

check_prerequisites() {
    if ! command -v sops &> /dev/null; then
        echo -e "${RED}Error: SOPS is not installed${NC}"
        echo "Install with: brew install sops"
        exit 1
    fi

    if ! command -v age &> /dev/null; then
        echo -e "${RED}Error: age is not installed${NC}"
        echo "Install with: brew install age"
        exit 1
    fi
}

get_key_file() {
    local env="$1"
    # Allow override via environment, else default
    if [[ -n "$SOPS_AGE_KEY_FILE" ]]; then
        echo "$SOPS_AGE_KEY_FILE"
    else
        echo "$HOME/.config/sops/age/keys.txt"
    fi
}

validate_environment() {
    local env="$1"
    if [[ ! "$env" =~ ^(local|staging|production)$ ]]; then
        echo -e "${RED}Error: Invalid environment '$env'${NC}"
        echo "Valid environments: local, staging, production"
        exit 1
    fi
}

validate_key_file() {
    local key_file="$(get_key_file)"
    
    if [[ ! -f "$key_file" ]]; then
        echo -e "${RED}Error: Age key file not found${NC}"
        echo "Expected: $key_file"
        echo ""
        echo "Generate keys with:"
        echo "  mkdir -p ~/.config/sops/age"
        echo "  age-keygen -o $key_file"
        exit 1
    fi
}

decrypt_file() {
    local env="$1"
    local file="$2"
    local key_file="$(get_key_file)"
    local secret_file="${SECRETS_DIR}/${env}/${file}"
    
    if [[ ! -f "$secret_file" ]]; then
        echo -e "${RED}Error: Secret file not found: $secret_file${NC}"
        exit 1
    fi
    
    echo -e "${BLUE}Decrypting $file for $env environment...${NC}"
    echo -e "Using key file: ${YELLOW}$key_file${NC}"
    SOPS_AGE_KEY_FILE="$key_file" sops -d -i "$secret_file"
    echo -e "${GREEN}✓ Decrypted successfully${NC}"
}

encrypt_file() {
    local env="$1"
    local file="$2"
    local key_file="$(get_key_file)"
    local secret_file="${SECRETS_DIR}/${env}/${file}"
    
    if [[ ! -f "$secret_file" ]]; then
        echo -e "${RED}Error: Secret file not found: $secret_file${NC}"
        exit 1
    fi
    
    echo -e "${BLUE}Encrypting $file for $env environment...${NC}"
    echo -e "Using key file: ${YELLOW}$key_file${NC}"
    SOPS_AGE_KEY_FILE="$key_file" sops -e -i "$secret_file"
    echo -e "${GREEN}✓ Encrypted successfully${NC}"
}

edit_file() {
    local env="$1"
    local file="$2"
    local key_file="$(get_key_file)"
    local secret_file="${SECRETS_DIR}/${env}/${file}"
    
    if [[ ! -f "$secret_file" ]]; then
        echo -e "${RED}Error: Secret file not found: $secret_file${NC}"
        exit 1
    fi
    
    echo -e "${BLUE}Editing $file for $env environment...${NC}"
    echo -e "Using key file: ${YELLOW}$key_file${NC}"
    SOPS_AGE_KEY_FILE="$key_file" sops "$secret_file"
    echo -e "${GREEN}✓ Edit completed${NC}"
}

view_file() {
    local env="$1"
    local file="$2"
    local key_file="$(get_key_file)"
    local secret_file="${SECRETS_DIR}/${env}/${file}"
    
    if [[ ! -f "$secret_file" ]]; then
        echo -e "${RED}Error: Secret file not found: $secret_file${NC}"
        exit 1
    fi
    
    echo -e "${BLUE}Viewing $file for $env environment:${NC}"
    echo ""
    echo -e "Using key file: ${YELLOW}$key_file${NC}"
    SOPS_AGE_KEY_FILE="$key_file" sops -d "$secret_file"
}

validate_all() {
    echo -e "${BLUE}Validating all encrypted secret files...${NC}"
    local errors=0
    local key_file="$(get_key_file)"
    
    if [[ ! -f "$key_file" ]]; then
        echo -e "${RED}  ✗ Age key file missing: $key_file${NC}"
        ((errors++))
        return 1
    fi
    
    for env in local staging production; do
        echo -e "${YELLOW}Checking $env environment:${NC}"
        
        if [[ ! -d "${SECRETS_DIR}/${env}" ]]; then
            echo -e "${RED}  ✗ Secrets directory missing: ${SECRETS_DIR}/${env}${NC}"
            ((errors++))
            continue
        fi
        
        for secret_file in "${SECRETS_DIR}/${env}"/*.yaml; do
            if [[ -f "$secret_file" ]]; then
                local filename=$(basename "$secret_file")
                if SOPS_AGE_KEY_FILE="$key_file" sops -d "$secret_file" >/dev/null 2>&1; then
                    echo -e "${GREEN}  ✓ $filename${NC}"
                else
                    echo -e "${RED}  ✗ $filename (decryption failed)${NC}"
                    ((errors++))
                fi
            fi
        done
    done
    
    if [[ $errors -eq 0 ]]; then
        echo -e "\n${GREEN}All secret files are valid!${NC}"
    else
        echo -e "\n${RED}Found $errors error(s)${NC}"
        exit 1
    fi
}

list_files() {
    local env="$1"
    echo -e "${BLUE}Secret files for $env environment:${NC}"
    echo ""
    
    if [[ ! -d "${SECRETS_DIR}/${env}" ]]; then
        echo -e "${RED}No secrets directory found for $env${NC}"
        exit 1
    fi
    
    for secret_file in "${SECRETS_DIR}/${env}"/*.yaml; do
        if [[ -f "$secret_file" ]]; then
            local filename=$(basename "$secret_file")
            echo "  $filename"
        fi
    done
}

generate_key() {
    local key_file="$(get_key_file)"
    
    if [[ -f "$key_file" ]]; then
        echo -e "${YELLOW}Warning: Age key file already exists: $key_file${NC}"
        read -p "Overwrite? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            echo -e "${BLUE}Cancelled${NC}"
            exit 0
        fi
    fi

    echo -e "${BLUE}Generating new age key...${NC}"
    
    if ! command -v age-keygen &> /dev/null; then
        echo -e "${RED}Error: age-keygen not found. Please install age: https://age-encryption.org/${NC}"
        exit 1
    fi

    mkdir -p "$(dirname "$key_file")"
    age-keygen -o "$key_file"
    chmod 600 "$key_file"
    
    echo -e "${GREEN}✓ Key generated: $key_file${NC}"
    echo -e "${BLUE}Public key:${NC}"
    grep "public key:" "$key_file"
    echo -e "${YELLOW}Remember to update .sops.yaml with the new public key for all environments!${NC}"
}

show_public_key() {
    local key_file="$(get_key_file)"
    
    if [[ ! -f "$key_file" ]]; then
        echo -e "${RED}Error: Age key file not found: $key_file${NC}"
        exit 1
    fi

    echo -e "${BLUE}Public key:${NC}"
    grep "public key:" "$key_file" | sed 's/# public key: //'
}

check_key_security() {
    echo -e "${BLUE}Checking key security...${NC}"
    local errors=0
    local key_file="$(get_key_file)"
    
    if [[ -f "$key_file" ]]; then
        local perms=$(stat -f "%A" "$key_file" 2>/dev/null || stat -c "%a" "$key_file" 2>/dev/null)
        if [[ "$perms" = "600" ]]; then
            echo -e "${GREEN}  ✓ Age key permissions: $perms (secure)${NC}"
        else
            echo -e "${YELLOW}  ⚠ Age key permissions: $perms (should be 600)${NC}"
            echo -e "    Fix with: chmod 600 $key_file"
            ((errors++))
        fi
    else
        echo -e "${RED}  ✗ Age key file not found: $key_file${NC}"
        ((errors++))
    fi
    
    # Check if keys are in gitignore
    if [[ -f "${PROJECT_ROOT}/.gitignore" ]] && grep -q "\.config/sops" "${PROJECT_ROOT}/.gitignore"; then
        echo -e "${GREEN}  ✓ SOPS config directory is gitignored${NC}"
    else
        echo -e "${YELLOW}  ⚠ Consider adding '.config/sops/' to .gitignore for security${NC}"
    fi
    
    if [[ $errors -eq 0 ]]; then
        echo -e "\n${GREEN}All security checks passed!${NC}"
    else
        echo -e "\n${YELLOW}Found $errors security issue(s)${NC}"
    fi
}

doctor() {
    local env="$1"
    local file="$2"
    local key_file="$(get_key_file)"
    local secret_file="${SECRETS_DIR}/${env}/${file}"

    echo -e "${BLUE}SOPS Doctor${NC}"
    echo -e "Environment: ${YELLOW}$env${NC}"
    echo -e "Secret file: ${YELLOW}$secret_file${NC}"
    echo -e "Key file:    ${YELLOW}$key_file${NC}"

    if [[ ! -f "$key_file" ]]; then
        echo -e "${RED}  ✗ Key file not found${NC}"
    else
        local my_pub
        my_pub=$(grep "public key:" "$key_file" | sed 's/# public key: //')
        if [[ -n "$my_pub" ]]; then
            echo -e "${GREEN}  ✓ Your public key: ${YELLOW}$my_pub${NC}"
        else
            echo -e "${YELLOW}  ⚠ Could not read public key from key file${NC}"
        fi
    fi

    # From .sops.yaml creation rules
    local cfg_pub
    cfg_pub=$(awk '/path_regex: deploy\/secrets\/'"$env"'\/.+\\.yaml\$/ {f=1} f && /age:/ {print $2; exit}' "$PROJECT_ROOT/.sops.yaml" 2>/dev/null || true)
    if [[ -n "$cfg_pub" ]]; then
        echo -e "${GREEN}  ✓ Configured recipient in .sops.yaml: ${YELLOW}$cfg_pub${NC}"
    else
        echo -e "${YELLOW}  ⚠ No matching creation_rule found in .sops.yaml for $env${NC}"
    fi

    # From the encrypted file metadata
    if [[ -f "$secret_file" ]]; then
        local file_recips
        file_recips=$(grep -E "^[[:space:]]*-?[[:space:]]*recipient:" "$secret_file" | awk '{print $2}' | paste -sd "," -)
        if [[ -n "$file_recips" ]]; then
            echo -e "${GREEN}  ✓ Recipients in file: ${YELLOW}$file_recips${NC}"
        else
            echo -e "${YELLOW}  ⚠ Could not detect recipients in file (file may not be SOPS-encrypted)${NC}"
        fi
    fi

    echo ""
    echo "Hints:"
    echo "- Ensure your key file contains the private key matching one of the recipients above."
    echo "- Override key path temporarily: SOPS_AGE_KEY_FILE=/path/to/keys.txt $0 edit $env $file"
    echo "- Show your public key to compare: $0 keyshow"
}

# Main script logic
check_prerequisites

case "${1:-}" in
    decrypt)
        if [[ $# -ne 3 ]]; then
            echo -e "${RED}Error: decrypt requires environment and file name${NC}"
            usage
            exit 1
        fi
        validate_environment "$2"
        validate_key_file
        decrypt_file "$2" "$3"
        ;;
    encrypt)
        if [[ $# -ne 3 ]]; then
            echo -e "${RED}Error: encrypt requires environment and file name${NC}"
            usage
            exit 1
        fi
        validate_environment "$2"
        validate_key_file
        encrypt_file "$2" "$3"
        ;;
    edit)
        if [[ $# -ne 3 ]]; then
            echo -e "${RED}Error: edit requires environment and file name${NC}"
            usage
            exit 1
        fi
        validate_environment "$2"
        validate_key_file
        edit_file "$2" "$3"
        ;;
    view)
        if [[ $# -ne 3 ]]; then
            echo -e "${RED}Error: view requires environment and file name${NC}"
            usage
            exit 1
        fi
        validate_environment "$2"
        validate_key_file
        view_file "$2" "$3"
        ;;
    validate)
        validate_all
        ;;
    list)
        if [[ $# -ne 2 ]]; then
            echo -e "${RED}Error: list requires environment${NC}"
            usage
            exit 1
        fi
        validate_environment "$2"
        list_files "$2"
        ;;
    keygen)
        if [[ $# -ne 1 ]]; then
            echo -e "${RED}Error: keygen doesn't require parameters${NC}"
            usage
            exit 1
        fi
        generate_key
        ;;
    keyshow)
        if [[ $# -ne 1 ]]; then
            echo -e "${RED}Error: keyshow doesn't require parameters${NC}"
            usage
            exit 1
        fi
        show_public_key
        ;;
    keycheck)
        check_key_security
        ;;
    updatekeys)
        if [[ $# -lt 2 || $# -gt 3 ]]; then
            echo -e "${RED}Error: updatekeys requires environment and optional file name${NC}"
            usage
            exit 1
        fi
        validate_environment "$2"
        env="$2"
        key_file="$(get_key_file)"
        if [[ -n "$3" ]]; then
            file="$3"
            secret_file="${SECRETS_DIR}/${env}/${file}"
            echo -e "${BLUE}Re-wrapping keys for ${YELLOW}$secret_file${NC} using ${YELLOW}$key_file${NC}..."
            SOPS_AGE_KEY_FILE="$key_file" sops updatekeys "$secret_file"
        else
            echo -e "${BLUE}Re-wrapping keys for all ${YELLOW}$env${NC} secrets using ${YELLOW}$key_file${NC}..."
            for secret_file in "${SECRETS_DIR}/${env}"/*.yaml; do
                [[ -f "$secret_file" ]] || continue
                echo "  - $(basename "$secret_file")"
                SOPS_AGE_KEY_FILE="$key_file" sops updatekeys "$secret_file"
            done
        fi
        ;;
    doctor)
        if [[ $# -ne 3 ]]; then
            echo -e "${RED}Error: doctor requires environment and file name${NC}"
            usage
            exit 1
        fi
        validate_environment "$2"
        doctor "$2" "$3"
        ;;
    *)
        usage
        exit 1
        ;;
esac
