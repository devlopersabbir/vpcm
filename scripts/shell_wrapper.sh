
# VPSM ssh wrapper override
ssh() {
    if [ $# -eq 0 ]; then
        command ssh
        return
    fi
    if [[ "$1" =~ ^[0-9]+$ ]] || [[ "$1" == *@* ]] || vpsm server list 2>/dev/null | grep -q -E "\b$1\b"; then
        vpsm ssh "$@"
    else
        command ssh "$@"
    fi
}
