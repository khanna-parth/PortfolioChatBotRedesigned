import { AppConfig } from "@/config/env";

export const HOST = `http://${AppConfig.backendUrl}`;
export const RAW_HOST = `${AppConfig.backendUrl}`

// export const HOST = 'http://localhost:3000';
// export const RAW_HOST = 'localhost:3000'
// export const HOST = 'http://169.254.86.5:3000';
// export const RAW_HOST = '169.254.86.5:3000';
// export const RAW_HOST = "api.chatbot.parthkhanna.me"
// export const HOST = `https://${RAW_HOST}`
export const HEARTBEAT_ROUTE = `${HOST}/heartbeat`
export const WEBSOCKET_ROUTE = `ws://${RAW_HOST}/ws`;
export const LIST_DOCS_ROUTE = `${HOST}/list-docs`;
export const DELETE_DOC_ROUTE = `${HOST}/delete-doc`
export const PROMPT_ROUTE = `${HOST}/prompt`;
export const UPLOAD_DOC_ROUTE = `${HOST}/add-doc`;
export const PRESETS_ROUTE = `${HOST}/presets`;
export const PRESET_PROMPT_ROUTE = `${HOST}/demo`
export const SUGGESTIONS_ROUTE = `${HOST}/demo/suggestions`