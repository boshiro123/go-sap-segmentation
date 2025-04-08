CREATE TABLE IF NOT EXISTS segmentation (
    id SERIAL PRIMARY KEY,
    address_sap_id VARCHAR(255) NOT NULL,
    adr_segment VARCHAR(16) NOT NULL,
    segment_id BIGINT NOT NULL,
    CONSTRAINT unique_address_sap_id UNIQUE (address_sap_id)
);

-- Создание индекса для быстрого поиска по address_sap_id
CREATE INDEX IF NOT EXISTS idx_segmentation_address_sap_id ON segmentation (address_sap_id);

-- Комментарии к таблице и столбцам
COMMENT ON TABLE segmentation IS 'Таблица для хранения данных сегментации из SAP';
COMMENT ON COLUMN segmentation.id IS 'Автоинкрементируемое уникальное поле';
COMMENT ON COLUMN segmentation.address_sap_id IS 'Идентификатор адреса в SAP';
COMMENT ON COLUMN segmentation.adr_segment IS 'Сегмент адреса';
COMMENT ON COLUMN segmentation.segment_id IS 'Идентификатор сегмента'; 