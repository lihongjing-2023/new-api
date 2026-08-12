/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

For commercial licensing, please contact support@quantumnous.com
*/
import { zodResolver } from '@hookform/resolvers/zod'
import { useForm, type Resolver } from 'react-hook-form'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'
import { z } from 'zod'

import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'

import { SettingsForm } from '../components/settings-form-layout'
import { SettingsPageFormActions } from '../components/settings-page-context'
import { SettingsSection } from '../components/settings-section'
import { useUpdateOption } from '../hooks/use-update-option'

const OPTION_KEY = 'traffic_fee_setting.traffic_fee'

const schema = z.object({
  trafficFee: z.coerce.number().min(0),
})

type Values = z.infer<typeof schema>

export function TrafficFeeSettingsSection({
  defaultValues,
}: {
  defaultValues: {
    trafficFee: number
  }
}) {
  const { t } = useTranslation()
  const updateOption = useUpdateOption()

  const form = useForm<Values>({
    resolver: zodResolver(schema) as unknown as Resolver<Values>,
    defaultValues: {
      trafficFee: defaultValues.trafficFee,
    },
  })

  const { isDirty, isSubmitting } = form.formState

  async function onSubmit(values: Values) {
    if (values.trafficFee === defaultValues.trafficFee) {
      toast.info(t('No changes to save'))
      return
    }

    await updateOption.mutateAsync({
      key: OPTION_KEY,
      value: values.trafficFee,
    })

    form.reset(values)
  }

  return (
    <SettingsSection title={t('Traffic Fee')}>
      <Form {...form}>
        <SettingsForm onSubmit={form.handleSubmit(onSubmit)} autoComplete='off'>
          <SettingsPageFormActions
            onSave={form.handleSubmit(onSubmit)}
            isSaving={updateOption.isPending || isSubmitting}
            isSaveDisabled={!isDirty}
            saveLabel='Save traffic fee'
          />
          <FormField
            control={form.control}
            name='trafficFee'
            render={({ field }) => (
              <FormItem>
                <FormLabel>{t('Traffic fee (USD per request)')}</FormLabel>
                <FormControl>
                  <Input
                    type='number'
                    min={0}
                    step={0.0001}
                    placeholder='0.001'
                    {...field}
                  />
                </FormControl>
                <FormDescription>
                  {t(
                    'Per-request surcharge applied at final billing for metered text models and async task settlement. Set to 0 to disable.'
                  )}
                </FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />
        </SettingsForm>
      </Form>
    </SettingsSection>
  )
}
