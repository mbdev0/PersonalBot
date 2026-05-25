export const PROGRAMS = [
  { label: 'Pumpfun Native', value: 'PumpfunNative' },
  { label: 'Pumpfun AMM', value: 'PumpfunAMM' },
] as const;

export type Program = (typeof PROGRAMS)[number]['value'];
