import assert from "node:assert/strict";
import * as StellarSdk from "@stellar/stellar-sdk";

const { rpc, xdr } = StellarSdk;
const keyXdr =
  "AAAABgAAAAGTnozAkG0Lx6x5RJuwKN73cnabLJM2wp+KFE2/pm6VqwAAABAAAAABAAAAAQAAAA4AAAAOc3RlbGxhci1yZXBsYXkAAAAAAAE=";

// This is the real API used by stellar/js-stellar-sdk's Soroban tests. The
// only change is that the endpoint is Stellar Replay's loopback server.
const server = new rpc.Server("http://127.0.0.1:8787", { allowHttp: true });
const ledgerKey = xdr.LedgerKey.fromXdr(keyXdr, "base64");
const response = await server.getLedgerEntries(ledgerKey);

assert.equal(response.latestLedger, 420);
assert.equal(response.entries.length, 1);
assert.equal(response.entries[0].lastModifiedLedgerSeq, 42);
assert.equal(response.entries[0].liveUntilLedgerSeq, 1000);
assert.equal(response.entries[0].key.toXdr("base64"), keyXdr);

console.log("stellar/js-stellar-sdk consumed Stellar Replay fixture successfully");
