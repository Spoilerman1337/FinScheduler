export type CalendarEventColorPreset =
    'red' | 'blue' | 'yellow' | 'green' | 'white' | 'orange' | 'violet';

export interface CalendarMarkerPreview {
    id: string;
    label: string;
    timeHints: string[];
    color: CalendarEventColorPreset;
    description: string;
    triggerTags: string[];
}

export interface CalendarDayDetails {
    dayKey: string;
    summary: string;
    note: string;
    markers: CalendarMarkerPreview[];
}

export interface EventCatalogEntry {
    id: string;
    dayKey: string;
    date: Date;
    dateLabel: string;
    fullDateLabel: string;
    summary: string;
    note: string;
    label: string;
    description: string;
    timeHints: string[];
    color: CalendarEventColorPreset;
    triggerTags: string[];
}

interface CalendarMarkerBlueprint extends Omit<CalendarMarkerPreview, 'id'> {
    id: string;
}

interface CalendarDayBlueprint {
    summary: string;
    note: string;
    markers: CalendarMarkerBlueprint[];
}

const monthFormatter = new Intl.DateTimeFormat('ru-RU', {
    month: 'long',
    year: 'numeric',
});

const monthShortFormatter = new Intl.DateTimeFormat('ru-RU', {
    month: 'short',
});

const dayMonthFormatter = new Intl.DateTimeFormat('ru-RU', {
    day: 'numeric',
    month: 'long',
});

const weekdayFormatter = new Intl.DateTimeFormat('ru-RU', {
    weekday: 'long',
});

function capitalizeLabel(value: string) {
    if (value.length === 0) {
        return value;
    }

    return value[0].toUpperCase() + value.slice(1);
}

function buildPreviewBlueprints(): CalendarDayBlueprint[] {
    return [
        {
            summary: 'Пакет автоплатежей',
            note: 'Здесь удобно держать регулярные списания и быстрые пометки по их статусу.',
            markers: [
                {
                    id: 'autopay-review',
                    label: 'Проверить автосписание',
                    timeHints: ['09:00'],
                    color: 'blue',
                    description: 'Сверить дату списания и убедиться, что сумма не изменилась.',
                    triggerTags: ['Автоплатеж', 'Контроль', 'Утро'],
                },
            ],
        },
        {
            summary: 'Подтвердить расходы',
            note: 'День можно выделить для ручной сверки трат перед финальной фиксацией.',
            markers: [
                {
                    id: 'receipts-check',
                    label: 'Сверка чеков',
                    timeHints: ['11:30'],
                    color: 'yellow',
                    description:
                        'Собрать чеки и отметить позиции, которые требуют дополнительной проверки.',
                    triggerTags: ['Чеки', 'Контроль', 'День'],
                },
                {
                    id: 'cashback-check',
                    label: 'Кешбэк',
                    timeHints: ['12:15'],
                    color: 'green',
                    description: 'Проверить, что кешбэк по выбранным операциям уже учтен.',
                    triggerTags: ['Кешбэк', 'Проверка', 'День'],
                },
            ],
        },
        {
            summary: 'Напоминание по налогам',
            note: 'Подходит для будущих обязательных платежей и задач, которые нельзя пропустить.',
            markers: [
                {
                    id: 'tax-reminder',
                    label: 'Налоговый дедлайн',
                    timeHints: ['14:00'],
                    color: 'yellow',
                    description:
                        'День заранее отмечен как чувствительный к срокам и подтверждениям.',
                    triggerTags: ['Налоги', 'Дедлайн', 'День'],
                },
            ],
        },
        {
            summary: 'День сверки бюджета',
            note: 'Подходит для сверки плана, факта и точек, где бюджет заметно отклонился от ожиданий.',
            markers: [
                {
                    id: 'budget-review',
                    label: 'Проверка категорий',
                    timeHints: ['10:00', '13:30', '18:00'],
                    color: 'blue',
                    description:
                        'Посмотреть крупные категории, где траты заметно выбились из плана.',
                    triggerTags: ['Бюджет', 'Контроль', 'Несколько слотов'],
                },
                {
                    id: 'limits-review',
                    label: 'Лимиты на неделю',
                    timeHints: [],
                    color: 'violet',
                    description:
                        'Подготовить лимиты на следующую неделю и зафиксировать, где нужен более жесткий контроль.',
                    triggerTags: ['Лимиты', 'План', 'Без времени'],
                },
            ],
        },
        {
            summary: 'Перевод между счетами',
            note: 'Подходит для дней, когда перевод нужно не только запланировать, но и сразу проконтролировать.',
            markers: [
                {
                    id: 'transfer-plan',
                    label: 'Запланировать перевод',
                    timeHints: ['13:00'],
                    color: 'green',
                    description:
                        'Выделить сумму и заранее отметить, какой счет будет источником.',
                    triggerTags: ['Перевод', 'План', 'День'],
                },
                {
                    id: 'transfer-confirmation',
                    label: 'Подтвердить зачисление',
                    timeHints: ['15:30'],
                    color: 'violet',
                    description:
                        'Проверить, что деньги пришли на нужный счет, и отметить перевод как завершенный.',
                    triggerTags: ['Перевод', 'Подтверждение', 'День'],
                },
            ],
        },
        {
            summary: 'Подписки и сервисы',
            note: 'Пометка для дат, где нужно не забыть про регулярные сервисы и их стоимость.',
            markers: [
                {
                    id: 'subscriptions',
                    label: 'Проверить подписки',
                    timeHints: ['18:00'],
                    color: 'blue',
                    description:
                        'Открыть список сервисов и решить, какие из них пора отключить или перенести.',
                    triggerTags: ['Подписки', 'Контроль', 'Вечер'],
                },
            ],
        },
    ];
}

function getDaysInMonth(date: Date) {
    return new Date(date.getFullYear(), date.getMonth() + 1, 0).getDate();
}

function isSameMonth(left: Date, right: Date) {
    return left.getFullYear() === right.getFullYear() && left.getMonth() === right.getMonth();
}

export function startOfDay(date: Date) {
    return new Date(date.getFullYear(), date.getMonth(), date.getDate());
}

export function startOfMonth(date: Date) {
    return new Date(date.getFullYear(), date.getMonth(), 1);
}

export function buildDayKey(date: Date) {
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');

    return `${year}-${month}-${day}`;
}

export function dayKeyToDate(dayKey: string) {
    const [year, month, day] = dayKey.split('-').map(Number);

    return new Date(year, month - 1, day);
}

export function formatMonthLabel(date: Date) {
    return capitalizeLabel(monthFormatter.format(startOfMonth(date)).replace(' г.', ''));
}

export function formatOverflowDayLabel(date: Date) {
    const shortMonth = monthShortFormatter.format(date).replace('.', '');

    return `${date.getDate()} ${shortMonth}`;
}

export function formatDayButtonLabel(date: Date) {
    return `${dayMonthFormatter.format(date)} ${date.getFullYear()}`;
}

export function formatSelectedDayLabel(date: Date) {
    return `${capitalizeLabel(weekdayFormatter.format(date))}, ${dayMonthFormatter.format(date)}`;
}

export function buildMonthPreviewDetails(displayedMonth: Date, referenceDate: Date = new Date()) {
    const monthStart = startOfMonth(displayedMonth);
    const daysInMonth = getDaysInMonth(monthStart);
    const previewDays = new Set<number>([3, 7, 12, 18, 24, Math.max(daysInMonth - 2, 1)]);

    if (isSameMonth(monthStart, referenceDate)) {
        previewDays.add(referenceDate.getDate());
    }

    const blueprints = buildPreviewBlueprints();
    const detailsByDay: Record<string, CalendarDayDetails> = {};
    const sortedDays = Array.from(previewDays)
        .filter((dayNumber) => dayNumber >= 1 && dayNumber <= daysInMonth)
        .sort((left, right) => left - right);

    sortedDays.forEach((dayNumber, index) => {
        const blueprint = blueprints[index % blueprints.length];
        const date = new Date(monthStart.getFullYear(), monthStart.getMonth(), dayNumber);
        const dayKey = buildDayKey(date);

        detailsByDay[dayKey] = {
            dayKey,
            summary: blueprint.summary,
            note: blueprint.note,
            markers: blueprint.markers.map((marker, markerIndex) => ({
                ...marker,
                id: `${dayKey}-${marker.id}-${markerIndex}`,
            })),
        };
    });

    return detailsByDay;
}

export function buildEventCatalogEntries(displayedMonth: Date, referenceDate: Date = new Date()) {
    const detailsByDay = buildMonthPreviewDetails(displayedMonth, referenceDate);

    return Object.values(detailsByDay)
        .sort((left, right) => left.dayKey.localeCompare(right.dayKey))
        .flatMap((details) => {
            const date = dayKeyToDate(details.dayKey);

            return details.markers.map<EventCatalogEntry>((marker) => ({
                id: marker.id,
                dayKey: details.dayKey,
                date,
                dateLabel: dayMonthFormatter.format(date),
                fullDateLabel: formatSelectedDayLabel(date),
                summary: details.summary,
                note: details.note,
                label: marker.label,
                description: marker.description,
                timeHints: marker.timeHints,
                color: marker.color,
                triggerTags: marker.triggerTags,
            }));
        });
}

export function getInitialSelectedDayKey(
    displayedMonth: Date,
    detailsByDay: Record<string, CalendarDayDetails>,
    referenceDate: Date = new Date(),
) {
    const monthStart = startOfMonth(displayedMonth);

    if (isSameMonth(monthStart, referenceDate)) {
        return buildDayKey(referenceDate);
    }

    const firstMarkedDay = Object.keys(detailsByDay).sort()[0];

    if (firstMarkedDay) {
        return firstMarkedDay;
    }

    return buildDayKey(monthStart);
}
