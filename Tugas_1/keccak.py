from eth_hash.auto import keccak
preimage = keccak.new(b'Hallo-1')
print(preimage.digest())
preimage_copy = preimage.copy()
preimage.update(b'hallo-2')
print(preimage.digest())
preimage_copy.update(b'hallo-3')
print(preimage_copy.digest())
