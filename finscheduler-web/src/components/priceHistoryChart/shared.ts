import type {LucideIcon} from 'lucide-react';

export const PRICE_HISTORY_TOOLTIP_WIDTH_PX = 184;
export const PRICE_HISTORY_ACTUAL_STROKE = 'var(--chakra-colors-neon-purple)';
export const PRICE_HISTORY_ACTUAL_FILL = 'var(--chakra-colors-neon-pink)';
export const PRICE_HISTORY_FORECAST_STROKE = 'var(--chakra-colors-neon-blue)';
export const PRICE_HISTORY_FORECAST_FILL = 'var(--chakra-colors-neon-cyan)';

export const PRICE_HISTORY_PRICE_INCREASE_COLOR = 'var(--chakra-colors-red-400)';
export const PRICE_HISTORY_PRICE_UNCHANGED_COLOR = 'var(--chakra-colors-neon-yellow)';
export const PRICE_HISTORY_PRICE_DECREASE_COLOR = 'var(--chakra-colors-neon-green)';

export interface PriceHistoryPoint {
    point: string;
    value: number;
    absoluteChange?: number | null;
    percentChange?: number | null;
}

export interface PriceHistoryChartPoint {
    point: string;
    actualValue: number | null;
    actualAbsoluteChange: number | null;
    actualPercentChange: number | null;
    forecastValue: number | null;
    forecastAbsoluteChange: number | null;
    forecastPercentChange: number | null;
}

export interface PriceHistoryChartMouseState {
    activeCoordinate?: {
        x?: number;
        y?: number;
    };
    activeTooltipIndex?: number | string | null;
    chartX?: number;
    chartY?: number;
}

export interface PriceHistoryLegendItem {
    label: string;
    color: string;
}

export interface PriceHistoryMetricItem {
    label: string;
    value: string;
    accentColor: string;
}

export interface PriceHistoryTooltipEntry {
    key: string;
    value: string;
    valueColor?: string;
    changeSummary?: string;
    changeColor?: string;
    changeIcon?: LucideIcon;
}

export interface PriceHistoryTooltipPosition {
    left: string;
    top: string;
    transform: string;
}

export interface PriceHistoryTooltipContent {
    label: string;
    entries: PriceHistoryTooltipEntry[];
    position: PriceHistoryTooltipPosition;
}
