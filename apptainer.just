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

shell CPUS="4" MEM="12G" NODE="" PARTITION="cpunormal":
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
        --cleanenv \
        --no-home \
        --bind $PWD:/workspace \
        --bind $HOME/.cache:$HOME/.cache \
        --bind $HOME/.claude:$HOME/.claude \
        --bind $HOME/.claude.json:$HOME/.claude.json \
        --bind $HOME/.codex:$HOME/.codex \
        --bind $HOME/.npm:$HOME/.npm \
        --bind $HOME/.config:$HOME/.config \
        --bind $HOME/.gitconfig:$HOME/.gitconfig:ro \
        --bind $HOME/.ssh:$HOME/.ssh:ro \
        --bind $HOME/.nvm:$HOME/.nvm \
        --bind $HOME/local/bin:$HOME/local/bin:ro \
        --env "BASH_ENV=/etc/profile.d/container-env.sh" \
        --env "GOCACHE=/tmp/go-build" \
        --env "GOMODCACHE=/tmp/go-mod" \
        --pwd /workspace \
        {{sif_image}}

