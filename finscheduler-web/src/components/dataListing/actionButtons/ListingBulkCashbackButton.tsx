import {Button} from '@chakra-ui/react';
import {PercentIcon} from 'lucide-react';

interface ListingBulkCashbackButtonProps {
    onClick: () => void;
}

export default function ListingBulkCashbackButton(props: ListingBulkCashbackButtonProps) {
    const {onClick} = props;

    return (
        <Button
            onClick={onClick}
            color="neon.blue"
            borderColor="neon.blue"
            backdropFilter="blur(12px)"
            bg="glass.bgHover"
            border="1px solid"
            transition="all 0.3s ease-in-out"
            _hover={{
                filter: 'drop-shadow(0 0 16px rgba(0, 212, 255, 0.7))',
                boxShadow: '0 0 20px rgba(0, 212, 255, 0.55)',
                color: 'neon.blue',
                bg: 'glass.bgHover',
                backdropFilter: 'blur(12px)',
                borderColor: 'neon.blue',
            }}
            focusRing="none"
        >
            <PercentIcon size={18} style={{marginRight: '8px'}} />
            Массовое обновление кешбека
        </Button>
    );
}
