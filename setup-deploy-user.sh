#!/bin/bash
# Setup Deploy User Script for Fedora VPS
# Run as root: sudo bash setup-deploy-user.sh

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}  Deploy User Setup Script (Fedora)    ${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""

# Check if running as root
if [ "$EUID" -ne 0 ]; then
    echo -e "${RED}Error: Please run as root (sudo bash setup-deploy-user.sh)${NC}"
    exit 1
fi

# Prompt for username
read -p "Enter deploy username (e.g., momentum-deploy): " USERNAME

if [ -z "$USERNAME" ]; then
    echo -e "${RED}Error: Username cannot be empty${NC}"
    exit 1
fi

# Prompt for project name (used for /opt directory)
read -p "Enter project name for /opt directory (e.g., momentum-business): " PROJECT_NAME

if [ -z "$PROJECT_NAME" ]; then
    echo -e "${RED}Error: Project name cannot be empty${NC}"
    exit 1
fi

echo ""
echo -e "${YELLOW}Creating user: ${USERNAME}${NC}"
echo -e "${YELLOW}Project directory: /opt/${PROJECT_NAME}${NC}"
echo ""
read -p "Continue? (y/n): " CONFIRM

if [ "$CONFIRM" != "y" ]; then
    echo "Aborted."
    exit 0
fi

echo ""

# Step 1: Create user
echo -e "${GREEN}[1/7] Creating user ${USERNAME}...${NC}"
if id "$USERNAME" &>/dev/null; then
    echo -e "${YELLOW}User ${USERNAME} already exists, skipping creation${NC}"
else
    useradd -r -s /bin/bash -m -d /home/${USERNAME} ${USERNAME}
    echo -e "User created"
fi

# Step 2: Add to docker group
echo -e "${GREEN}[2/7] Adding ${USERNAME} to docker group...${NC}"
usermod -aG docker ${USERNAME}
echo -e "Added to docker group"

# Step 3: Create project directory
echo -e "${GREEN}[3/7] Creating /opt/${PROJECT_NAME}...${NC}"
mkdir -p /opt/${PROJECT_NAME}
chown ${USERNAME}:${USERNAME} /opt/${PROJECT_NAME}
echo -e "Directory created and owned by ${USERNAME}"

# Step 4: Create SSH directory and generate key
echo -e "${GREEN}[4/7] Setting up SSH keys...${NC}"
mkdir -p /home/${USERNAME}/.ssh
ssh-keygen -t ed25519 -f /home/${USERNAME}/.ssh/github_deploy -N "" -C "${USERNAME}@github-actions"
echo -e "SSH keypair generated"

# Step 5: Setup authorized_keys
echo -e "${GREEN}[5/7] Configuring authorized_keys...${NC}"
cat /home/${USERNAME}/.ssh/github_deploy.pub > /home/${USERNAME}/.ssh/authorized_keys
echo -e "Public key added to authorized_keys"

# Step 6: Fix permissions
echo -e "${GREEN}[6/7] Setting permissions...${NC}"
chmod 700 /home/${USERNAME}/.ssh
chmod 600 /home/${USERNAME}/.ssh/authorized_keys
chmod 600 /home/${USERNAME}/.ssh/github_deploy
chmod 644 /home/${USERNAME}/.ssh/github_deploy.pub
chown -R ${USERNAME}:${USERNAME} /home/${USERNAME}/.ssh
echo -e "Permissions set"

# Step 7: SELinux context (Fedora)
echo -e "${GREEN}[7/7] Restoring SELinux context...${NC}"
if command -v restorecon &> /dev/null; then
    restorecon -R /home/${USERNAME}/.ssh
    echo -e "SELinux context restored"
else
    echo -e "${YELLOW}restorecon not found, skipping SELinux${NC}"
fi

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}  Setup Complete!                      ${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo -e "${YELLOW}Summary:${NC}"
echo -e "  User:              ${USERNAME}"
echo -e "  Home:              /home/${USERNAME}"
echo -e "  Project Dir:       /opt/${PROJECT_NAME}"
echo -e "  Docker Access:     Yes"
echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}  PRIVATE KEY (for GitHub VPS_SSH_KEY) ${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
cat /home/${USERNAME}/.ssh/github_deploy
echo ""
echo -e "${GREEN}========================================${NC}"
echo ""
echo -e "${YELLOW}GitHub Actions Secrets to set:${NC}"
echo -e "  VPS_HOST:     $(hostname -I | awk '{print $1}')"
echo -e "  VPS_USER:     ${USERNAME}"
echo -e "  VPS_SSH_KEY:  (copy the private key above)"
echo ""
echo -e "${YELLOW}To save the key on your local machine:${NC}"
echo -e "  1. Create file: ~/.ssh/${USERNAME}"
echo -e "  2. Paste the private key above"
echo -e "  3. Run: chmod 600 ~/.ssh/${USERNAME}"
echo ""
echo -e "${YELLOW}To test SSH connection:${NC}"
echo -e "  ssh -i ~/.ssh/${USERNAME} ${USERNAME}@$(hostname -I | awk '{print $1}')"
echo ""
