# Apptainer Go development container management

# Set the SIF image name
sif_image := "dev/bookworm.sif"

# List all apptainer recipes
default:
    @just --list apptainer

# Build the container from Apptainer definition file
build:
    #!/usr/bin/env bash
    set -euo pipefail
    mkdir -p dev
    if ! command -v apptainer &> /dev/null; then
        module load apptainer/1.4.1-nkna
    fi
    apptainer build {{sif_image}} dev/bookworm.def

# Rebuild the container from scratch (forces rebuild)
rebuild:
    #!/usr/bin/env bash
    set -euo pipefail
    rm -f {{sif_image}}
    mkdir -p dev
    if ! command -v apptainer &> /dev/null; then
        module load apptainer/1.4.1-nkna
    fi
    apptainer build {{sif_image}} dev/bookworm.def

# Drop into a bash shell in the container (local)
shell:
    #!/usr/bin/env bash
    set -euo pipefail
    if ! command -v apptainer &> /dev/null; then
        module load apptainer/1.4.1-nkna
    fi
    apptainer shell \
        --no-home \
        --bind $PWD:/workspace \
        --bind $HOME/.claude:$HOME/.claude \
        --bind $HOME/.claude.json:$HOME/.claude.json \
        --bind $HOME/.codex:$HOME/.codex \
        --bind $HOME/.config:$HOME/.config \
        --bind $HOME/.gitconfig:$HOME/.gitconfig:ro \
        --bind $HOME/.ssh:$HOME/.ssh:ro \
        --bind $HOME/.bashrc:$HOME/.bashrc:ro \
        --bind $HOME/.nvm:$HOME/.nvm \
        --bind $HOME/local/bin:$HOME/local/bin:ro \
        --env "PATH=$HOME/local/bin:$PATH" \
        --pwd /workspace \
        {{sif_image}}

# Interactive shell using SLURM (8 CPUs, 32GB RAM, with home mount)
# Optionally specify CPUS, MEM, NODE, and PARTITION
shell-slurm CPUS="4" MEM="12G" NODE="" PARTITION="cpunormal":
    #!/usr/bin/env bash
    set -euo pipefail
    if ! command -v apptainer &> /dev/null; then
        module load apptainer/1.4.1-nkna
    fi

    SRUN_CMD="srun --cpus-per-task={{CPUS}} --mem={{MEM}} --pty"
    if [ -n "{{PARTITION}}" ]; then
        SRUN_CMD="$SRUN_CMD --partition={{PARTITION}}"
    fi
    if [ -n "{{NODE}}" ]; then
        SRUN_CMD="$SRUN_CMD --nodelist={{NODE}}"
        echo "Requesting node {{NODE}} with {{CPUS}} CPUs, {{MEM}} RAM..."
    else
        echo "Requesting {{CPUS}} CPUs, {{MEM}} RAM (auto-selecting node)..."
    fi

    $SRUN_CMD apptainer shell \
        --no-home \
        --bind $PWD:/workspace \
        --bind $HOME/.claude:$HOME/.claude \
        --bind $HOME/.claude.json:$HOME/.claude.json \
        --bind $HOME/.codex:$HOME/.codex \
        --bind $HOME/.config:$HOME/.config \
        --bind $HOME/.gitconfig:$HOME/.gitconfig:ro \
        --bind $HOME/.ssh:$HOME/.ssh:ro \
        --bind $HOME/.bashrc:$HOME/.bashrc:ro \
        --bind $HOME/.nvm:$HOME/.nvm \
        --bind $HOME/local/bin:$HOME/local/bin:ro \
        --env "PATH=$HOME/local/bin:$PATH" \
        --pwd /workspace \
        {{sif_image}}

# Run a command in the container (local)
run CMD:
    #!/usr/bin/env bash
    set -euo pipefail
    if ! command -v apptainer &> /dev/null; then
        module load apptainer/1.4.1-nkna
    fi
    apptainer exec \
        --no-home \
        --bind $PWD:/workspace \
        --bind $HOME/.claude:$HOME/.claude \
        --bind $HOME/.claude.json:$HOME/.claude.json \
        --bind $HOME/.codex:$HOME/.codex \
        --bind $HOME/.config:$HOME/.config \
        --bind $HOME/.gitconfig:$HOME/.gitconfig:ro \
        --bind $HOME/.ssh:$HOME/.ssh:ro \
        --bind $HOME/.bashrc:$HOME/.bashrc:ro \
        --bind $HOME/.nvm:$HOME/.nvm \
        --bind $HOME/local/bin:$HOME/local/bin:ro \
        --env "PATH=$HOME/local/bin:$PATH" \
        --pwd /workspace \
        {{sif_image}} {{CMD}}

# Run command using SLURM
run-slurm CMD CPUS="4" MEM="12G" NODE="" PARTITION="cpunormal":
    #!/usr/bin/env bash
    set -euo pipefail
    if ! command -v apptainer &> /dev/null; then
        module load apptainer/1.4.1-nkna
    fi

    SRUN_CMD="srun --cpus-per-task={{CPUS}} --mem={{MEM}} 
    if [ -n "{{PARTITION}}" ]; then
        SRUN_CMD="$SRUN_CMD --partition={{PARTITION}}"
    fi
    if [ -n "{{NODE}}" ]; then
        SRUN_CMD="$SRUN_CMD --nodelist={{NODE}}"
    fi

    $SRUN_CMD apptainer exec \
        --no-home \
        --bind $PWD:/workspace \
        --bind $HOME/.claude:$HOME/.claude \
        --bind $HOME/.claude.json:$HOME/.claude.json \
        --bind $HOME/.codex:$HOME/.codex \
        --bind $HOME/.config:$HOME/.config \
        --bind $HOME/.gitconfig:$HOME/.gitconfig:ro \
        --bind $HOME/.ssh:$HOME/.ssh:ro \
        --bind $HOME/.bashrc:$HOME/.bashrc:ro \
        --bind $HOME/.nvm:$HOME/.nvm \
        --bind $HOME/local/bin:$HOME/local/bin:ro \
        --env "PATH=$HOME/local/bin:$PATH" \
        --pwd /workspace \
        {{sif_image}} {{CMD}}

# Test the container by running version checks
test:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "Testing container tools..."
    if ! command -v apptainer &> /dev/null; then
        module load apptainer/1.4.1-nkna
    fi
    apptainer exec \
        --no-home \
        --bind $PWD:/workspace \
        --bind $HOME/.claude:$HOME/.claude \
        --bind $HOME/.claude.json:$HOME/.claude.json \
        --bind $HOME/.codex:$HOME/.codex \
        --bind $HOME/.config:$HOME/.config \
        --bind $HOME/.gitconfig:$HOME/.gitconfig:ro \
        --bind $HOME/.ssh:$HOME/.ssh:ro \
        --bind $HOME/.bashrc:$HOME/.bashrc:ro \
        --bind $HOME/.nvm:$HOME/.nvm \
        --bind $HOME/local/bin:$HOME/local/bin:ro \
        --env "PATH=$HOME/local/bin:$PATH" \
        --pwd /workspace \
        {{sif_image}} bash -c "\
        echo 'Go:' && go version && \
        echo 'Node.js:' && node --version && \
        echo 'npm:' && npm --version && \
        echo 'Git:' && git --version && \
        echo 'Make:' && make --version | head -1 && \
        echo 'Bash:' && bash --version | head -1 && \
        echo 'Claude Code:' && npm list -g @anthropic-ai/claude-code --depth=0 | grep claude-code \
    "

# Show container info
info:
    #!/usr/bin/env bash
    set -euo pipefail
    if [ -f {{sif_image}} ]; then
        if ! command -v apptainer &> /dev/null; then
            module load apptainer/1.4.1-nkna
        fi
        apptainer inspect {{sif_image}}
    else
        echo "Container image {{sif_image}} not found. Run 'just apptainer build' first."
    fi

# Clean up the container image
clean:
    rm -f {{sif_image}}
    @echo "Removed {{sif_image}}"
