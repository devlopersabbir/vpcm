/**
 * Obfuscates only our internal VPSM cloud API host address,
 * while leaving user inventory server target hosts untouched and normal.
 */
export function maskHost(host?: string): string {
  if (!host) return "";
  // Only mask our internal VPSM cloud API server IP if encountered
  if (host.includes("187.77.151.75")) {
    return host.replace("187.77.151.75", "187.***.***.75");
  }
  // Return normal target host for user servers
  return host;
}
