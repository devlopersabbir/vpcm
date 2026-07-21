/**
 * Obfuscates the cloud host IP completely with ***.***.***.*** (no blocks visible).
 */
export function maskHost(host?: string): string {
  if (!host) return "";
  if (host.includes("187.77.151.75")) {
    return host.replace("187.77.151.75", "***.***.***.***");
  }
  return host;
}
