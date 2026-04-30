import { Card } from '@/components/ui/card';
import type { Settings } from '../types/settings';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Pencil, Save } from 'lucide-react';
import { usePostSettings } from '../hooks/useSettings';
import { toast } from 'sonner';

interface QuickSellButtonsProps {
  data: Settings;
}

export function QuickSellButtons({ data }: QuickSellButtonsProps) {
  const [isEditable, setIsEditable] = useState(false);
  const [button1, setButton1] = useState(data.quick_sell_buttons.button_1 * 100);
  const [button2, setButton2] = useState(data.quick_sell_buttons.button_2 * 100);
  const [button3, setButton3] = useState(data.quick_sell_buttons.button_3 * 100);
  const [button4, setButton4] = useState(data.quick_sell_buttons.button_4 * 100);

  const postSettingsMutation = usePostSettings();

  const update = (patch: Partial<Settings>) => {
    if (!data) return;
    postSettingsMutation.mutate(
      { ...data, ...patch },
      {
        onSuccess: () => {
          toast.success('Succesfully updated Quick Sell Buttons');
        },
        onError: (error) => {
          toast.error('Failed to update Quick Sell Buttons', { description: error.message });
        },
      }
    );
  };

  return (
    <Card className="flex flex-col flex-1 rounded-lg border border-foreground/10 bg-foreground/3 px-4">
      <p className="text-sm font-semibold text-foreground/50 uppercase">Quick Sell Buttons</p>

      <div className="grid grid-rows-2 grid-cols-4 gap-2">
        <Label className="justify-self-center">Button 1</Label>
        <Label className="justify-self-center">Button 2</Label>
        <Label className="justify-self-center">Button 3</Label>
        <Label className="justify-self-center">Button 4 </Label>

        <QSButtonInput
          isEditable={isEditable}
          buttonVal={button1}
          onChange={setButton1}
        ></QSButtonInput>

        <QSButtonInput
          isEditable={isEditable}
          buttonVal={button2}
          onChange={setButton2}
        ></QSButtonInput>

        <QSButtonInput
          isEditable={isEditable}
          buttonVal={button3}
          onChange={setButton3}
        ></QSButtonInput>

        <QSButtonInput
          isEditable={isEditable}
          buttonVal={button4}
          onChange={setButton4}
        ></QSButtonInput>
      </div>

      <div className="flex justify-end gap-1.5 mt-auto">
        <Button
          className="hover:bg-yellow-300 hover:opacity-50"
          onClick={() => setIsEditable(!isEditable)}
        >
          <Pencil />
        </Button>
        <Button
          className="hover:bg-green-400 hover:opacity-50"
          onClick={() =>
            update({
              quick_sell_buttons: {
                button_1: button1 / 100,
                button_2: button2 / 100,
                button_3: button3 / 100,
                button_4: button4 / 100,
              },
            })
          }
        >
          <Save />
        </Button>
      </div>
    </Card>
  );
}

interface QSButtonInputProps {
  isEditable: boolean;
  buttonVal: number;
  onChange: (n: number) => void;
}

function QSButtonInput({ isEditable, buttonVal, onChange }: QSButtonInputProps) {
  return (
    <Input
      className=""
      type="number"
      max={100}
      min={0}
      disabled={!isEditable}
      value={buttonVal}
      onChange={(e) => {
        let number = e.target.valueAsNumber;

        if (number < 0) {
          number = buttonVal;
        }

        if (number > 100) {
          number = buttonVal;
        }

        onChange(number);
      }}
    ></Input>
  );
}
