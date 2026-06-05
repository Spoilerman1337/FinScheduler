export function getCashbackColor(cashback: number | undefined): string {
    let color = 'neon.pink';

    if (cashback) {
        color = cashback > 1 ? 'neon.green' : 'neon.yellow';
    }

    return color;
}
