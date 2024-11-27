genesis_dirs=("val-alice" "val-bob" "node-carol")
for dir in "${genesis_dirs[@]}"; do
    echo "Removing directory: $dir"
    rm -rf "$HOME/mandu-devnet/$dir"
done
