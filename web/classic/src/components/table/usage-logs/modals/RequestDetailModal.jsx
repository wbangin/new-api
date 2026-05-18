import React from 'react';
import { Modal, Spin, Typography, Tabs, TabPane } from '@douyinfe/semi-ui';
import { useTranslation } from 'react-i18next';

const { Text } = Typography;

function formatJson(str) {
  if (!str) return '';
  try {
    const obj = JSON.parse(str);
    return JSON.stringify(obj, null, 2);
  } catch (e) {
    return str;
  }
}

const RequestDetailModal = ({
  showRequestDetailModal,
  setShowRequestDetailModal,
  requestDetailData,
  requestDetailLoading,
}) => {
  const { t } = useTranslation();

  return (
    <Modal
      title={t('请求/响应详情')}
      visible={showRequestDetailModal}
      onCancel={() => setShowRequestDetailModal(false)}
      footer={null}
      width={800}
      style={{ maxHeight: '80vh' }}
      bodyStyle={{ maxHeight: '65vh', overflow: 'auto' }}
    >
      {requestDetailLoading ? (
        <div style={{ textAlign: 'center', padding: 40 }}>
          <Spin size='large' />
        </div>
      ) : requestDetailData ? (
        <div>
          <div style={{ marginBottom: 12, display: 'flex', gap: 16, flexWrap: 'wrap' }}>
            <Text type='tertiary'>{t('模型')}: {requestDetailData.model}</Text>
            <Text type='tertiary'>{t('令牌')}: {requestDetailData.token_name}</Text>
            <Text type='tertiary'>{t('状态码')}: {requestDetailData.status_code}</Text>
            <Text type='tertiary'>{t('流式')}: {requestDetailData.is_stream ? t('是') : t('否')}</Text>
          </div>
          <Tabs type='line'>
            <TabPane tab={t('请求体')} itemKey='request'>
              <pre style={{
                background: '#f5f5f5',
                padding: 12,
                borderRadius: 8,
                overflow: 'auto',
                maxHeight: '45vh',
                fontSize: 12,
                lineHeight: 1.5,
                whiteSpace: 'pre-wrap',
                wordBreak: 'break-all',
              }}>
                {formatJson(requestDetailData.request_body) || t('无数据')}
              </pre>
            </TabPane>
            <TabPane tab={t('响应体')} itemKey='response'>
              <pre style={{
                background: '#f5f5f5',
                padding: 12,
                borderRadius: 8,
                overflow: 'auto',
                maxHeight: '45vh',
                fontSize: 12,
                lineHeight: 1.5,
                whiteSpace: 'pre-wrap',
                wordBreak: 'break-all',
              }}>
                {formatJson(requestDetailData.response_body) || t('无数据（流式响应不记录响应体）')}
              </pre>
            </TabPane>
          </Tabs>
        </div>
      ) : (
        <div style={{ textAlign: 'center', padding: 40 }}>
          <Text type='tertiary'>{t('未找到请求详情数据')}</Text>
        </div>
      )}
    </Modal>
  );
};

export default RequestDetailModal;
