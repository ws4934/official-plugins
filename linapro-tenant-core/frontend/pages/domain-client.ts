import { pluginApiPath, requestClient } from "#/api/request";

const pluginID = "linapro-tenant-core";

export type DomainStatus = "active" | "disabled";

export interface PlatformDomain {
  id: number;
  tenantId: number;
  domain: string;
  isPrimary: boolean;
  isVerified: boolean;
  status: DomainStatus;
  createdAt?: number | null;
}

export interface PlatformDomainListParams {
  pageNum?: number;
  pageSize?: number;
  tenantId?: number;
  domain?: string;
  status?: DomainStatus | "";
}

export interface PlatformDomainPayload {
  tenantId: number;
  domain: string;
  isPrimary?: boolean;
}

export async function platformDomainList(params?: PlatformDomainListParams) {
  const res = await requestClient.get<{
    list: PlatformDomain[];
    total: number;
  }>(pluginApiPath(pluginID, "platform/domains"), { params });
  return { items: res.list, total: res.total };
}

export function platformDomainCreate(payload: PlatformDomainPayload) {
  return requestClient.post<PlatformDomain>(
    pluginApiPath(pluginID, "platform/domains"),
    payload,
  );
}

export function platformDomainSetVerified(id: number, verified: boolean) {
  return requestClient.put(
    pluginApiPath(pluginID, `platform/domains/${id}/verification`),
    { verified },
  );
}

export function platformDomainDelete(id: number) {
  return requestClient.delete(
    pluginApiPath(pluginID, `platform/domains/${id}`),
  );
}
