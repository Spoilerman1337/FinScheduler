import type {CalendarEventColorPreset} from './shared.ts';

interface CalendarEventColorStyle {
    bg: string;
    glow: string;
    text: string;
}

const calendarEventColorPresets = [
    'red',
    'blue',
    'yellow',
    'green',
    'white',
    'orange',
    'violet',
] as const satisfies readonly CalendarEventColorPreset[];

const calendarEventColorStyles: Record<CalendarEventColorPreset, CalendarEventColorStyle> = {
    red: {
        bg: 'app.negative',
        glow: 'none',
        text: 'app.negative',
    },
    blue: {
        bg: 'app.info',
        glow: 'app.glowCyan',
        text: 'app.info',
    },
    yellow: {
        bg: 'amber.300',
        glow: 'none',
        text: 'amber.300',
    },
    green: {
        bg: 'app.positive',
        glow: 'none',
        text: 'app.positive',
    },
    white: {
        bg: 'fg',
        glow: 'none',
        text: 'fg',
    },
    orange: {
        bg: 'amber.500',
        glow: 'none',
        text: 'amber.500',
    },
    violet: {
        bg: 'app.accentViolet',
        glow: 'app.glowViolet',
        text: 'app.accentViolet',
    },
};

export {calendarEventColorPresets, calendarEventColorStyles};
