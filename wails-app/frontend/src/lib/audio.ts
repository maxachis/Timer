let audioCtx: AudioContext | null = null;
function getAudioCtx(): AudioContext {
  return (audioCtx ||= new AudioContext());
}

export function playAlertSound() {
  const ctx = getAudioCtx();
  const osc = ctx.createOscillator();
  const gain = ctx.createGain();
  osc.connect(gain);
  gain.connect(ctx.destination);
  osc.type = "sine";
  osc.frequency.setValueAtTime(880, ctx.currentTime);
  osc.frequency.setValueAtTime(660, ctx.currentTime + 0.15);
  osc.frequency.setValueAtTime(880, ctx.currentTime + 0.3);
  gain.gain.setValueAtTime(0.3, ctx.currentTime);
  gain.gain.exponentialRampToValueAtTime(0.01, ctx.currentTime + 0.5);
  osc.start(ctx.currentTime);
  osc.stop(ctx.currentTime + 0.5);
}

export function playBlip(freq: number) {
  const ctx = getAudioCtx();
  const t0 = ctx.currentTime;
  const osc = ctx.createOscillator();
  const gain = ctx.createGain();
  osc.type = "sine";
  osc.frequency.value = freq;
  osc.connect(gain).connect(ctx.destination);
  gain.gain.setValueAtTime(0, t0);
  gain.gain.linearRampToValueAtTime(0.2, t0 + 0.005);
  gain.gain.exponentialRampToValueAtTime(0.0001, t0 + 0.13);
  osc.start(t0);
  osc.stop(t0 + 0.15);
}

export const playStartBlip = () => playBlip(880);
export const playStopBlip = () => playBlip(440);
