import base58
import base64
from borsh_construct import CStruct, String, U64, I64, Bool, Bytes
import json


public_key = Bytes[16]

args_schema = CStruct(
    "name" / String,
    "symbol" / String,
    "uri" / String
)

sell_schema = CStruct(
    
)

# The base64 data string from the transaction
# base58_data = "58yiJoQi6w3VrqFxNVubx8s9YM2UvLw5qvUeinV6LWFXXykLdF7UV66QUf8sSZN4BHHbPrEHPb8T1kxXsNaGvWoztbmtiMmCvzLvEz59LjTifqv1a9WA8KkKejA8bDyAUbjki2X5D6"
# base58_data = "U5L9Srg16xJaQ1bJQiNLih4DTf8XESsggzLChuEWTn64tbNDp4GwPkYhuLFKUF536xLmQZiPC8xK8xmqoq4BJajmabsZh4Ed88FTC6dbYuZPMvK5GvU9ogQCwxeqZHwrHPBgoQ39g4vRurT"
# base58_data = "h95wiy926iNkGitSi3kh1S3v7Qci9BcGQZMRghq5NJAuxBm9QPjpLPHbkpWi38TMR2JghSeWVQ4zU8PwrYuBKyin6NgbmNnHqvWfB5BAq7acfvnyv5xAd782CEpcuENX8zLSPDDQGVyjZP3L4Xx"

# \x18\x1e\xc8(\x05\x1c\x07w\x05\x00\x00\x00Fizzi\x05\x00\x00\x00FIZZIG\x00\x00\x00https://cf-ipfs.com/ipfs/QmcjpoXkVfzqwmVGusCsqAFiigVCKgvCtFLLvZTrUkKx5m'
# first 6 bytes are SOMETHING (dint know what)
# next byte is length of name (0x05 = 5) in hex
# x00 is padding
# next 5 bytes are name
# next byte is length of symbol (0x05 = 5) in hex
# x00 is padding
# next 5 bytes are symbol
# then more padding -> url 

""""fields": [
        {
          "name": "mint",
          "type": "publicKey",
          "index": false
        },
        {
          "name": "solAmount",
          "type": "u64",
          "index": false
        },
        {
          "name": "tokenAmount",
          "type": "u64",
          "index": false
        },
        {
          "name": "isBuy",
          "type": "bool",
          "index": false
        },
        {
          "name": "user",
          "type": "publicKey",
          "index": false
        },
        {
          "name": "timestamp",
          "type": "i64",
          "index": false
        },
        {
          "name": "virtualSolReserves",
          "type": "u64",
          "index": false
        },
        {
          "name": "virtualTokenReserves",
          "type": "u64",
          "index": false
        }
      ]
      
    """

trade_event = CStruct(
    "mint" / Bytes[32],
    "solAmount" / U64,
    "tokenAmount" / U64,
    "isBuy" / Bool,
    "user" / public_key,
    "timestamp" / I64,
    "virtualSolReserves" / U64,
    "virtualTokenReserves" / U64
)

create_even = CStruct(
    "name" / String,
    "symbol" / String,
    "uri" / String,
    "publicKey" / public_key,
    "bondingCurve" / public_key,
    "user" / public_key
)

sell_event = CStruct(
    "amount" / U64,
    "minSolOutput" / U64,
)

base58_d = "58yiJoQi6w3VrqFxNVubx8s9YM2UvLw5qvUeinV6LWFXXykLdF7UV66QUf8sSZN4BHHbPrEHPb8T1kxXsNaGvWoztbmtiMmCvzLvEz59LjTifqv1a9WA8KkKejA8bDyAUbjki2X5D6"
decoded = base58.b58decode(base58_d)
print(list(decoded))