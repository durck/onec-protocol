# Security-provider integration

Evidence: E12. The 1C-specific contract is Name(), Step(target,input) returning
output bytes/local completion, and cleanup. Tokens are opaque to the 1C framing
layer. See [bootstrap](03-bootstrap.md) and [inner HTTP security](07-native-authentication.md)
for placement. A provider context is never shared between manager, worker and
inner exchanges, even when they use one credential identity.

## Windows Negotiate

The live baseline obtains outbound credentials for package Negotiate and uses
the provider's credential-name attribute in V8IntroTicketReq. It invokes
InitializeSecurityContext with the acknowledgment target, TOKEN buffers,
request flags `0x113` and native data representation `0x10`. The reference adds
OS-managed output allocation. Copy output bytes before releasing the OS buffer.

Handle SEC_E_OK, SEC_I_CONTINUE_NEEDED, SEC_I_COMPLETE_NEEDED and
SEC_I_COMPLETE_AND_CONTINUE. The COMPLETE statuses require CompleteAuthToken.
Only OK or COMPLETE_NEEDED means local completion. Pass continuation/final input
back to the same context. After completion, query the negotiated package rather
than equating the offered name Negotiate with Kerberos. Delete the context and
free the credential handle separately. The OS API defines these calls and
statuses: [InitializeSecurityContext](https://learn.microsoft.com/en-us/windows/win32/secauthn/initializesecuritycontext--negotiate).

This uses the OS security implementation, not any 1C library. Successful provider
authentication neither proves infobase login nor enables automatic signing or
encryption of later application records.

## Explicit NTLMv2 profile

This independent provider has synthetic interoperability evidence, not live
domain-backed 1C acceptance. It emits raw NTLMSSP type 1/2/3 messages to the
Negotiate acceptor, without a SPNEGO wrapper. Introduction identity is
`EXAMPLE\alice` for the synthetic domain/user. The credentials are an OS account
password or its NT hash, not an infobase SecPswd value.

The standard message layouts and algorithms are defined by
[MS-NLMP](https://learn.microsoft.com/en-us/openspecs/windows_protocols/ms-nlmp/b38c36ed-2804-4868-a9ff-8dd3182128e4).
The retained profile makes these choices:

- Require UNICODE, NTLM, EXTENDED_SESSIONSECURITY, TARGET_INFO and 128-bit flags;
  also offer REQUEST_TARGET, SIGN, ALWAYS_SIGN and KEY_EXCH.
- Do not advertise VERSION. The initial message is 32 bytes, with empty domain
  and workstation buffers. Keep the exact type-1/type-2 bytes for MIC.
- Derive NT hash as MD4(UTF16LE(password)); response key is
  HMAC-MD5(NThash, UTF16LE(Uppercase(user) || domain)). Preserve domain case.
- Generate a fresh eight-byte client challenge per context. Use the server's
  timestamp if present, otherwise current FILETIME. Use a 24-byte zero LM field.
- Construct the NTLMv2 blob with response version 1/1, zero reserved fields,
  timestamp, client challenge, response AV pairs and final four zero bytes.
  NT proof is HMAC-MD5(responseKey, serverChallenge || blob); NT response is
  proof || blob; session base key is HMAC-MD5(responseKey, proof).
- With SIGN and KEY_EXCH, generate a fresh 16-byte exported session key and
  RC4-encrypt it using the session base key. Otherwise use the session base key.
- Type 3 reserves an all-zero eight-byte Version area and a 16-byte MIC at
  offset 72; variable payload starts at 88. MIC is HMAC-MD5(exportedKey,
  type1 || type2 || type3WithZeroMIC). Response AV flags indicate MIC.
- No workstation name or local OS version is published. Client target-name AV
  is empty; channel-binding AV is zero. Reject nonzero server channel binding.

Validate signatures, types, negotiated flags, security-buffer bounds/overlap,
UTF-16 and terminated AV pairs. The profile rejects duplicate AVs, datagram,
anonymous, LM-key, local-call and unrequested sealing modes. It has no LM/NTLMv1
fallback. Final type 3 establishes local completion; an unexpected final server
token is an error. Enclosing 1C completion must still succeed.

NTLM does not supply Kerberos-style mutual peer proof. This profile does not
implement general application-message signing/sealing or TLS channel binding.
Servers requiring different wrapping or protection may reject it.

## Kerberos profile

The independent provider has synthetic evidence. Standard ticket/authenticator
semantics come from [RFC 4120](https://www.rfc-editor.org/rfc/rfc4120); GSS token
and checksum details come from [RFC 4121](https://www.rfc-editor.org/rfc/rfc4121).
The 1C integration uses fresh AP-REQs and validates final AP-REPs. Do not replay
captured authenticators.

Select a currently valid service ticket and associated session key for the
explicit configured SPN. A service ticket cannot be repurposed by renaming it.
The peer-announced target is not sufficient to choose an independently trusted
SPN. A same-realm TGT may obtain one service ticket through a separately
configured KDC; cross-realm referrals/renewal are outside the retained profile.

The constructed authenticator has checksum type `0x8003`: u32LE(16), sixteen
zero channel-binding bytes and u32LE flags `0x3e` (no delegation). APOptions
sets mutual-required. GSS framing is DER APPLICATION 0, containing the Kerberos
OID `1.2.840.113554.1.2.2`, token ID bytes `01 00`, then DER AP-REQ.

The accepted final reply has that same GSS OID and token ID `02 00`, then
AP-REP. Validate version 5, expected encryption type, decrypt with the selected
session key and key usage 12, and require exact echoed authenticator ctime/cusec.
Require complete DER consumption. Then report local completion with no outgoing
token. Do not let the server's 1C completion flag bypass this proof.

The retained input adapters accept unencrypted KRB-CRED containers and MIT FILE
ccache v3/v4 with one identity, associated keys and valid time ranges. They do
not import anything into OS caches. Credential-file formats are local adapters,
not part of the 1C wire protocol. Exact configured identity must remain stable
between outer and inner providers.

## Wrapping and acceptance boundaries

Negotiate, raw NTLMSSP, Kerberos GSS and SPNEGO are different concepts.
[RFC 4178](https://www.rfc-editor.org/rfc/rfc4178) defines SPNEGO; do not add or
strip a wrapper simply because 1C reports Negotiate. The raw profiles above
describe the implemented choices and their evidence, not every server policy.
Use standard provider libraries where possible and implement only the 1C
introduction/token routing here. Anonymous enumeration has its own bounded
acceptance and does not establish anonymous application access.
