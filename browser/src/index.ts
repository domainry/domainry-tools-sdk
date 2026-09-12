export interface ToolSetting {
 key: string; enabled: boolean; revision: number; updated_at?: string;
 version: string; description: string; available: boolean;
 state: 'available' | 'disabled' | 'connection_unavailable' | 'connection_unknown';
}
export interface ToolSettingInput { enabled: boolean; expected_revision: number; tool_version: string }
export interface ToolsClientDependencies {
 request<T>(path: string, options?: {method?: string; body?: unknown; signal?: AbortSignal}): Promise<T>;
}
// Owner-defined paths; the product supplies an authenticated same-origin transport.
export class ToolsClient {
 constructor(private readonly dependencies: ToolsClientDependencies) {}
 listToolSettings(signal?: AbortSignal) {
  return this.dependencies.request<{items: ToolSetting[]}>('/tools/preferences', {signal});
 }
 updateToolSetting(key: string, input: ToolSettingInput, signal?: AbortSignal) {
  if (!/^[a-zA-Z0-9_][a-zA-Z0-9_.-]*$/.test(key) || key === '.' || key === '..') throw new Error('tools.settings.key_invalid');
  return this.dependencies.request<ToolSetting>(`/tools/preferences/${encodeURIComponent(key)}`, {method:'PUT',body:input,signal});
 }
}
