import {screen, waitFor} from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import {useState} from 'react';
import {describe, expect, it, vi} from 'vitest';
import {renderWithProviders} from '../../test/render.tsx';
import AsyncSelectField from './AsyncSelectField.tsx';

function SingleSelectHarness() {
    const [value, setValue] = useState('');

    return (
        <AsyncSelectField
            label="Tag"
            value={value}
            placeholder="Select tag"
            loadOptions={vi.fn().mockResolvedValue({
                options: [{label: "Solomon's", value: 'tag-1'}],
                hasMore: false,
            })}
            onChange={setValue}
        />
    );
}

describe('AsyncSelectField', () => {
    it('keeps the search input empty after selecting an option in single select mode', async () => {
        // Arrange
        const user = userEvent.setup();

        renderWithProviders(<SingleSelectHarness />);

        // Act
        const combobox = screen.getByRole('combobox', {name: 'Tag'});
        await user.click(combobox);
        await user.click(await screen.findByRole('option', {name: "Solomon's"}));

        // Assert
        await waitFor(() => {
            expect(screen.queryByRole('listbox')).not.toBeInTheDocument();
        });
        expect(combobox).toHaveValue('');
    });
});
