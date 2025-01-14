#!/bin/bash

# Function to ensure bin directory exists
ensure_bin_dir() {
    if [ ! -d "bin" ]; then
        echo "Creating bin directory..."
        mkdir -p bin
    fi
}

# Function to compile specific example
compile_example() {
    local tag=$1
    local output="bin/$2"
    
    echo "Compiling $tag example..."
    go build -tags="$tag" -o "$output"
    
    if [ $? -eq 0 ]; then
        echo "Successfully compiled $output"
    else
        echo "Failed to compile $tag example"
        exit 1
    fi
}

# Function to show usage
show_usage() {
    echo "Usage: $0 [command]"
    echo "Commands:"
    echo "  (no args)   - compile both examples"
    echo "  precompiled - compile precompiled example only"
    echo "  basic       - compile basic example only"
    echo "  clean       - remove compiled binaries"
    exit 1
}

# Main compilation logic
if [ $# -eq 0 ]; then
    # No arguments - compile both examples
    ensure_bin_dir
    compile_example "precompiled" "precompiledExample"
    echo "-------------------"
    compile_example "basic" "basicExample"
elif [ $# -eq 1 ]; then
    # One argument - compile specific example or clean
    case $1 in
        "precompiled")
            ensure_bin_dir
            compile_example "precompiled" "precompiledExample"
            ;;
        "basic")
            ensure_bin_dir
            compile_example "basic" "basicExample"
            ;;
        "clean")
            echo "Cleaning compiled binaries..."
            rm -rf bin
            echo "Clean completed"
            ;;
        *)
            show_usage
            ;;
    esac
else
    show_usage
fi
