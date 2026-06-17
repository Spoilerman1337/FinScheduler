import {Box} from '@chakra-ui/react';
import {
    Area,
    AreaChart,
    Brush,
    CartesianGrid,
    Line,
    ResponsiveContainer,
    Tooltip,
    XAxis,
    YAxis,
} from 'recharts';
import {
    PRICE_HISTORY_ACTUAL_FILL,
    PRICE_HISTORY_ACTUAL_STROKE,
    PRICE_HISTORY_BRUSH_FILL,
    PRICE_HISTORY_BRUSH_STROKE,
    PRICE_HISTORY_FORECAST_FILL,
    PRICE_HISTORY_FORECAST_STROKE,
} from '../shared.ts';
import type {
    PriceHistoryChartMouseState,
    PriceHistoryChartPoint,
    PriceHistoryTooltipContent,
} from '../shared.ts';
import PriceHistoryTooltip from './PriceHistoryTooltip.tsx';

interface PriceHistoryChartCanvasProps {
    chartLabel: string;
    data: PriceHistoryChartPoint[];
    hasForecast: boolean;
    activeTooltip: PriceHistoryTooltipContent | null;
    onMouseMove: (state: PriceHistoryChartMouseState) => void;
    onMouseLeave: () => void;
    gradientId: string;
    forecastGradientId: string;
    formatAxisDate: (value: string) => string;
    formatAxisMoney: (value: number) => string;
}

export default function PriceHistoryChartCanvas({
    chartLabel,
    data,
    hasForecast,
    activeTooltip,
    onMouseMove,
    onMouseLeave,
    gradientId,
    forecastGradientId,
    formatAxisDate,
    formatAxisMoney,
}: PriceHistoryChartCanvasProps) {
    const hasBrush = data.length > 2;

    return (
        <Box
            aria-label={chartLabel}
            position="relative"
            height={hasBrush ? '280px' : '240px'}
            borderRadius="xl"
            border="1px solid"
            borderColor="app.cardBorder"
            bg="linear-gradient(180deg, rgba(16, 29, 49, 0.92), rgba(8, 16, 32, 0.82))"
            px={{base: 2, md: 3}}
            py={3}
        >
            {activeTooltip ? (
                <Box
                    data-testid="price-history-tooltip"
                    position="absolute"
                    left={activeTooltip.position.left}
                    top={activeTooltip.position.top}
                    transform={activeTooltip.position.transform}
                    zIndex={2}
                    pointerEvents="none"
                >
                    <PriceHistoryTooltip
                        label={activeTooltip.label}
                        entries={activeTooltip.entries}
                    />
                </Box>
            ) : null}
            <ResponsiveContainer width="100%" height="100%" minWidth={320} minHeight={220}>
                <AreaChart
                    data={data}
                    margin={{top: 8, right: 4, left: 0, bottom: 0}}
                    onMouseMove={onMouseMove}
                    onMouseLeave={onMouseLeave}
                >
                    <defs>
                        <linearGradient id={gradientId} x1="0" y1="0" x2="0" y2="1">
                            <stop
                                offset="0%"
                                stopColor={PRICE_HISTORY_ACTUAL_FILL}
                                stopOpacity={0.3}
                            />
                            <stop
                                offset="100%"
                                stopColor={PRICE_HISTORY_ACTUAL_FILL}
                                stopOpacity={0.02}
                            />
                        </linearGradient>
                        <linearGradient id={forecastGradientId} x1="0" y1="0" x2="0" y2="1">
                            <stop
                                offset="0%"
                                stopColor={PRICE_HISTORY_FORECAST_FILL}
                                stopOpacity={0.22}
                            />
                            <stop
                                offset="100%"
                                stopColor={PRICE_HISTORY_FORECAST_FILL}
                                stopOpacity={0.01}
                            />
                        </linearGradient>
                    </defs>
                    <CartesianGrid
                        vertical={false}
                        stroke="rgba(183, 171, 255, 0.14)"
                        strokeDasharray="4 8"
                    />
                    <XAxis
                        dataKey="point"
                        axisLine={false}
                        tickLine={false}
                        minTickGap={24}
                        stroke="rgba(186, 199, 223, 0.72)"
                        tickFormatter={formatAxisDate}
                        style={{fontSize: '12px'}}
                    />
                    <YAxis
                        axisLine={false}
                        tickLine={false}
                        width={54}
                        stroke="rgba(186, 199, 223, 0.72)"
                        tickFormatter={formatAxisMoney}
                        style={{fontSize: '11px'}}
                    />
                    <Tooltip
                        content={() => null}
                        cursor={{
                            stroke: PRICE_HISTORY_ACTUAL_STROKE,
                            strokeOpacity: 0.28,
                        }}
                        wrapperStyle={{display: 'none'}}
                    />
                    <Area
                        type="monotone"
                        dataKey="actualValue"
                        stroke="none"
                        fill={`url(#${gradientId})`}
                        fillOpacity={1}
                        connectNulls
                    />
                    <Line
                        type="monotone"
                        dataKey="actualValue"
                        name="Факт"
                        stroke={PRICE_HISTORY_ACTUAL_STROKE}
                        strokeWidth={3}
                        dot={{
                            r: 4,
                            strokeWidth: 2,
                            fill: PRICE_HISTORY_ACTUAL_STROKE,
                            stroke: '#050816',
                        }}
                        activeDot={{
                            r: 6,
                            strokeWidth: 2,
                            fill: PRICE_HISTORY_ACTUAL_STROKE,
                            stroke: '#050816',
                        }}
                        connectNulls
                    />
                    {hasForecast ? (
                        <>
                            <Area
                                type="monotone"
                                dataKey="forecastValue"
                                stroke="none"
                                fill={`url(#${forecastGradientId})`}
                                fillOpacity={1}
                                connectNulls
                            />
                            <Line
                                type="monotone"
                                dataKey="forecastValue"
                                name="Прогноз"
                                stroke={PRICE_HISTORY_FORECAST_STROKE}
                                strokeWidth={3}
                                strokeDasharray="7 6"
                                dot={{
                                    r: 3.5,
                                    strokeWidth: 2,
                                    fill: PRICE_HISTORY_FORECAST_STROKE,
                                    stroke: '#050816',
                                }}
                                activeDot={{
                                    r: 5.5,
                                    strokeWidth: 2,
                                    fill: PRICE_HISTORY_FORECAST_STROKE,
                                    stroke: '#050816',
                                }}
                                connectNulls
                            />
                        </>
                    ) : null}
                    {hasBrush ? (
                        <Brush
                            ariaLabel="Диапазон графика"
                            dataKey="point"
                            height={28}
                            travellerWidth={10}
                            stroke={PRICE_HISTORY_BRUSH_STROKE}
                            fill={PRICE_HISTORY_BRUSH_FILL}
                            tickFormatter={formatAxisDate}
                        />
                    ) : null}
                </AreaChart>
            </ResponsiveContainer>
        </Box>
    );
}
