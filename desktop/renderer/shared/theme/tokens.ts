export interface ThemeTokens {
  bg: string; surface: string; text: string; muted: string; border: string;
  primary: string; secondary: string; success: string; warn: string; error: string;
  graphNode: string; graphEdge: string; graphSelected: string; graphHover: string;
}

export const requiredThemeKeys: Array<keyof ThemeTokens> = [
  'bg','surface','text','muted','border','primary','secondary','success','warn','error','graphNode','graphEdge','graphSelected','graphHover'
];
