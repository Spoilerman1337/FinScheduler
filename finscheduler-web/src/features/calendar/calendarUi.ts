import type {CalendarMarkerTone} from './shared.ts';

const markerToneStyles: Record<CalendarMarkerTone, {bg: string; text: string}> = {
    accent: {
        bg: 'app.accent',
        text: 'app.accent',
    },
    positive: {
        bg: 'app.positive',
        text: 'app.positive',
    },
    warning: {
        bg: 'app.warning',
        text: 'app.warning',
    },
    violet: {
        bg: 'app.accentViolet',
        text: 'app.accentViolet',
    },
};

export {markerToneStyles};
