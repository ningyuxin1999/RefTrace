import re
import sys

def transform_index_md(file_path):
    with open(file_path, 'r') as f:
        content = f.read()

    # Find the table section
    table_match = re.search(r'\|\s*\[`.*?\].*?\|\n\|[-\s|]*\n(\|.*?\n)*', content)
    if not table_match:
        return content

    table = table_match.group(0)
    
    # Split into rows
    rows = table.split('\n')
    # Keep header row and content rows, skip the separator
    kept_rows = [row for row in rows if row.strip() and not row.startswith('|----')]
    
    # Create new table with our header and all original rows
    new_table = "| Class | Description |\n|------|-------------|\n"
    new_table += '\n'.join(kept_rows)

    # Replace old table with new table
    new_content = content.replace(table, new_table)
    
    with open(file_path, 'w') as f:
        f.write(new_content)

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: transform_docs.py <index.md path>")
        sys.exit(1)
    transform_index_md(sys.argv[1])