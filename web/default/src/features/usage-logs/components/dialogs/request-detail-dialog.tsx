import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Copy, Check, Loader2 } from 'lucide-react'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Label } from '@/components/ui/label'
import { useCopyToClipboard } from '@/hooks/use-copy-to-clipboard'
import { getRequestDetail, type RequestDetailEntry } from '../../api'

interface RequestDetailDialogProps {
  requestId: string
  open: boolean
  onOpenChange: (open: boolean) => void
}

export function RequestDetailDialog(props: RequestDetailDialogProps) {
  const { t } = useTranslation()
  const { copiedText, copyToClipboard } = useCopyToClipboard({ notify: false })
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [data, setData] = useState<RequestDetailEntry | null>(null)

  const handleOpen = async (open: boolean) => {
    props.onOpenChange(open)
    if (open && !data && !loading) {
      setLoading(true)
      setError(null)
      try {
        const res = await getRequestDetail(props.requestId)
        if (res.success && res.data) {
          setData(res.data)
        } else {
          setError(res.message || t('Failed to load request detail'))
        }
      } catch (e: any) {
        setError(e?.message || t('Failed to load request detail'))
      } finally {
        setLoading(false)
      }
    }
  }

  const formatJson = (str: string): string => {
    try {
      return JSON.stringify(JSON.parse(str), null, 2)
    } catch {
      return str
    }
  }

  return (
    <Dialog open={props.open} onOpenChange={handleOpen}>
      <DialogContent className='max-w-3xl max-sm:max-w-[calc(100vw-1.5rem)] max-sm:p-4'>
        <DialogHeader>
          <DialogTitle>{t('Request & Response Detail')}</DialogTitle>
          <DialogDescription className='sr-only'>
            {t('View the complete request and response content')}
          </DialogDescription>
        </DialogHeader>

        <ScrollArea className='max-h-[70vh] pr-4 max-sm:max-h-[calc(100dvh-7rem)]'>
          {loading && (
            <div className='flex items-center justify-center py-12'>
              <Loader2 className='text-muted-foreground size-6 animate-spin' />
            </div>
          )}

          {error && (
            <div className='py-8 text-center text-sm text-red-500'>
              {error}
            </div>
          )}

          {data && (
            <div className='space-y-4 py-1'>
              {/* Request Body */}
              <div className='space-y-1.5'>
                <div className='flex items-center justify-between'>
                  <Label className='text-sm font-semibold'>
                    {t('Request Body')}
                  </Label>
                  <Button
                    variant='ghost'
                    size='sm'
                    className='h-6 px-2'
                    onClick={() => copyToClipboard(data.request_body)}
                  >
                    {copiedText === data.request_body ? (
                      <Check className='mr-1 size-3 text-green-600' />
                    ) : (
                      <Copy className='mr-1 size-3' />
                    )}
                    {t('Copy')}
                  </Button>
                </div>
                <pre className='bg-muted/50 max-h-80 overflow-auto rounded-md border p-3 font-mono text-xs leading-relaxed whitespace-pre-wrap'>
                  {formatJson(data.request_body) || t('(empty)')}
                </pre>
              </div>

              {/* Response Body */}
              <div className='space-y-1.5'>
                <div className='flex items-center justify-between'>
                  <Label className='text-sm font-semibold'>
                    {t('Response Body')}
                  </Label>
                  <Button
                    variant='ghost'
                    size='sm'
                    className='h-6 px-2'
                    onClick={() => copyToClipboard(data.response_body)}
                  >
                    {copiedText === data.response_body ? (
                      <Check className='mr-1 size-3 text-green-600' />
                    ) : (
                      <Copy className='mr-1 size-3' />
                    )}
                    {t('Copy')}
                  </Button>
                </div>
                <pre className='bg-muted/50 max-h-80 overflow-auto rounded-md border p-3 font-mono text-xs leading-relaxed whitespace-pre-wrap'>
                  {formatJson(data.response_body) || t('(empty)')}
                </pre>
              </div>

              {/* Metadata */}
              <div className='text-muted-foreground space-y-0.5 text-xs'>
                <div>
                  {t('Model')}: {data.model} | {t('Stream')}:{' '}
                  {data.is_stream ? t('Yes') : t('No')} | {t('Status')}:{' '}
                  {data.status_code}
                </div>
                <div>
                  {t('Time')}:{' '}
                  {new Date(data.timestamp * 1000).toLocaleString()}
                </div>
              </div>
            </div>
          )}
        </ScrollArea>
      </DialogContent>
    </Dialog>
  )
}
