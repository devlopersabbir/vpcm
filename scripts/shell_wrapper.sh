
# VPSM ssh wrapper override
ssh() {
    if [ $# -eq 0 ]; then
        command ssh
        return
    fi
    local use_vpsm=0
    for arg in "$@"; do
        if [[ "$arg" == -* ]]; then
            continue
        fi
        if [[ "$arg" =~ ^[0-9]+$ ]] || [[ "$arg" == *@* ]] || vpsm server list 2>/dev/null | grep -q -E "\b$arg\b"; then
            use_vpsm=1
            break
        fi
    done

    if [ $use_vpsm -eq 1 ]; then
        vpsm ssh "$@"
    else
        command ssh "$@"
    fi
}
