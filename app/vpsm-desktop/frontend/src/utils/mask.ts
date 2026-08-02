/**
 * Obfuscates the cloud host IP completely with ***.***.***.*** (no blocks visible).
 */
export function maskHost(host?: string): string {
  if (!host) return "";
  // Obfuscate standard IPv4 addresses if needed for privacy display
  return host.replace(/\b(?:\d{1,3}\.){3}\d{1,3}\b/g, "***.***.***.***");
}
