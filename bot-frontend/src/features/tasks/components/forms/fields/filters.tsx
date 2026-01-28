import { Input } from '@/components/ui/input';
import { Switch } from '@/components/ui/switch';
import type { Filters } from '@/features/tasks/types/strategyTask';
import { Label } from '@radix-ui/react-label';

interface FiltersEntryProps {
  filters: Filters;
  onFiltersChange: (filters: Filters) => void;
}

export function FiltersEntry({ filters, onFiltersChange }: FiltersEntryProps) {
  const updateFilter = (key: keyof Filters, checked: boolean) => {
    onFiltersChange({
      ...filters,
      [key]: checked ? true : null,
    });
  };

  const updateInputFilter = (key: keyof Filters, input: string | number) => {
    onFiltersChange({
      ...filters,
      [key]: input,
    });
  };
  return (
    <div className="flex flex-row justify-evenly gap-5 p-5">
      <div className="flex flex-col items-center space-y-2">
        <Label htmlFor="hasWebsite">Has Website</Label>
        <Switch
          id="has-website"
          checked={filters.has_website ?? false}
          onCheckedChange={(checked) => updateFilter('has_website', checked)}
        ></Switch>
      </div>

      <div className="flex flex-col items-center space-y-2">
        <Label htmlFor="hasTelegram">Has Telegram</Label>
        <Switch
          id="has-telegram"
          checked={filters.has_telegram ?? false}
          onCheckedChange={(checked) => updateFilter('has_telegram', checked)}
        ></Switch>
      </div>

      <div className="flex flex-col items-center space-y-2">
        <Label htmlFor="hasTwitter">Has Twitter</Label>
        <Switch
          id="has-twitter"
          checked={filters.has_twitter ?? false}
          onCheckedChange={(checked) => updateFilter('has_twitter', checked)}
        ></Switch>
      </div>

      <div className="flex flex-col items-center">
        <Label htmlFor="devWallet">Developer Wallet</Label>
        <Input
          id="devWallet"
          value={filters.dev_wallet ?? ''}
          placeholder="Developer Wallet"
          onChange={(e) => updateInputFilter('dev_wallet', e.target.value)}
        ></Input>
      </div>
    </div>
  );
}
