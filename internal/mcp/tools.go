package mcp

const (
	ToolResolveCustomer     = "resolve_customer"
	ToolResolveWarehouse    = "resolve_warehouse"
	ToolResolveProduct      = "resolve_product"
	ToolResolveMaterial     = "resolve_material"
	ToolProductDetails      = "product_details"
	ToolSalesReport         = "sales_report"
	ToolStockBalance        = "stock_balance"
	ToolAvailabilityReport  = "availability_report"
	ToolTopProducts         = "top_products"
	ToolCustomerSummary     = "customer_summary"
	ToolResolveSalesChannel = "resolve_sales_channel"
	ToolResolveFirm         = "resolve_firm"
	ToolResolveCash         = "resolve_cash"
	ToolResolveCostArticle  = "resolve_cost_article"
	ToolResolveOperation    = "resolve_operation"
	ToolCashBalance         = "cash_balance"
	ToolCashFlow            = "cash_flow"
	ToolReceivablesBalance  = "receivables_balance"
	ToolPayablesBalance     = "payables_balance"
	ToolPurchasesReport     = "purchases_report"
	ToolGoodsInTransit      = "goods_in_transit"
	ToolStockReserves       = "stock_reserves"
	ToolReturnsReport       = "returns_report"
	ToolEventLog            = "event_log"
	ToolObjectHistory       = "object_history"
	ToolFindDocument        = "find_document"

	// Производственный блок: спецификации (нормы расхода) и документы Производство.
	ToolProductSpecification     = "product_specification"
	ToolSpecificationCost        = "specification_cost"
	ToolSpecificationExplode     = "specification_explode"
	ToolSpecificationWhereUsed   = "specification_where_used"
	ToolSpecificationVersions    = "specification_versions"
	ToolSpecificationList        = "specification_list"
	ToolProductionOutput         = "production_output"
	ToolProductionConsumption    = "production_consumption"
	ToolProductionDocumentDetail = "production_document"
)

// ScopeReportCost — доступ к себестоимости. Исторически это measure-level право (меры
// cost/profit/margin внутри sales_report), поэтому оно проверяется отдельно от ToolScopes
// при фильтрации схемы в handleToolsList. Начиная с производственного блока это ещё и
// tool-level scope: нормы расхода сырья и себестоимость выпуска — та же чувствительная
// информация, что и закупочные цены, поэтому производственные инструменты закрыты им же
// (см. ToolScopes ниже). Финально право проверяется на стороне 1С по заголовку
// X-MCP-Scopes (defense in depth).
const ScopeReportCost = "mcp:report:cost"

// CostMeasures — меры, закрытые правом ScopeReportCost, по инструментам. Должны быть синхронны
// с белыми списками мер в CommonModules/MCP (BSL) и со значениями enum в схемах ниже.
//
// Ключ — инструмент, а не просто имя меры: имя `amount` само по себе о чувствительности ничего
// не говорит. В sales_report это выручка, в purchases_report — закупочные обороты, открытые
// money-правом. Закрывать нужно ровно те меры, за которыми стоит СЕБЕСТОИМОСТЬ, а это свойство
// пары «инструмент + мера», а не имени.
var CostMeasures = map[string][]string{
	ToolSalesReport: {"cost", "profit", "margin"},
	// stock_balance.amount — СуммаBalance регистра остатков, то есть себестоимость партий
	// (ЦенаПартии × Количество): amount/qty даёт закупочную цену за единицу.
	ToolStockBalance: {"amount"},
	// Товары в пути: обе суммы — стоимость поставки от поставщика, та же закупочная цена.
	// Без этого её читали бы отсюда в обход закрытого stock_balance.amount.
	ToolGoodsInTransit: {"amount", "amount_in_currency"},
}

// facetNotes — что сказать модели о грани, добавленной профилем базы (extra). Профиль
// передаёт только имена: смысл грани гейт знает, а 1С в health его не описывает.
type facetNotes struct {
	// Filters — описание свойства отбора целиком.
	Filters map[string]string
	// GroupBy — предложение, дописываемое к описанию group_by, если значение добавлено.
	GroupBy map[string]string
}

// extraFacetNotes — описания граней extra по инструментам. Грань без записи здесь всё равно
// добавляется, просто с общим описанием.
var extraFacetNotes = map[string]facetNotes{
	// Связь возврата с продажей (УПП): заказ из шапки возврата и документ продажи из его строк.
	ToolReturnsReport: {
		Filters: map[string]string{
			"order_ids": "Filter by customer order IDs (the order the return was entered against). " +
				"Take the id from find_document or from the order dimension of this report.",
			"sale_document_ids": "Filter by sale document IDs (the sale the returned goods came from). " +
				"Take the id from find_document (e.g. doc_type РеализацияТоваровУслуг) or from the " +
				"sale_document dimension. Returns are selected by THEIR date: to answer 'were there " +
				"returns for this sale', the period must cover the return dates, not the sale date.",
		},
		GroupBy: map[string]string{
			"order": "order = the customer order the return was entered against; the only link to " +
				"the sale for returns not based on a sale document.",
			"sale_document": "sale_document = the sale the returned goods came from; filled only for " +
				"returns entered on the basis of a sale document. Returns against an order (web " +
				"orders), manual returns and retail returns in the same shift show '(не указано)' — " +
				"for them use order.",
		},
	},
}
