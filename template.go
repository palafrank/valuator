package valuator

const valuatorTemplate = `
<html>
  <body>
  <h4>Valuation Metrics</h4>
    <table border="1">
      <tr>
        <th>
          Filed
        </th>
        <th>
          Book
        </th>
        <th>
          DPS
        </th>
        <th>
          FCF
        </th>
        <th>
          Payout
        </th>
        <th>
          NetMargin
        </th>
        <th>
          OpsMargin
        </th>
        <th>
          OpsLev
        </th>
        <th>
          FinLev
        </th>
        <th>
          CRatio
        </th>
        <th>
          RoA
        </th>
        <th>
          RoE
        </th>
      </tr>
      {{ range $index, $m := .FiledData }}
      <tr>
        <th>
          {{ $m.Date.String }}
        </th>
        <th>
          {{ $m.BookValue }}
        </th>
        <th>
          {{ $m.DividendPerShare }}
        </th>
        <th>
          {{ printf "%.2f" $m.FreeCashFlow }}
        </th>
        <th>
          {{ $m.PayOutToFcf }}
        </th>
        <th>
          {{ $m.ContribMargin }}
        </th>
        <th>
          {{ $m.OpsMargin }}
        </th>
        <th>
          {{ $m.OperatingLeverage }}
        </th>
        <th>
          {{ $m.FinancialLeverage }}
        </th>
        <th>
          {{ $m.CurrentRatio }}
        </th>
        <th>
          {{ $m.ReturnOnAssets }}
        </th>
        <th>
          {{ $m.ReturnOnEquity }}
        </th>
      </tr>
      {{ end }}
    </table>
    <h4>YoY Metrics</h4>
    <table border="1">
      <tr>
        <th>
          Filed
          </th>
        <th>
          Revenue(%)
        </th>
        <th>
          Earnings(%)
        </th>
        <th>
          FCF(%)
        </th>
        <th>
          Margin(%)
        </th>
        <th>
          Debt(%)
        </th>
        <th>
          Equity(%)
        </th>
        <th>
          BV
        </th>
        <th>
          Div
        </th>
      </tr>
      {{ range $index, $m := .FiledData }}
      {{ if (isYoyNonNil $m) }}
      {{ $yoy := $m.Yoy }}
      <tr>
        <th>
          {{ $m.Date.String }}
        </th>
        <th>
          {{ $yoy.RevenueGrowth }}
        </th>
        <th>
          {{ $yoy.EarningsGrowth }}
        </th>
        <th>
          {{ $yoy.CashFlowGrowth }}
        </th>
        <th>
          {{ $yoy.GrossMarginGrowth }}
        </th>
        <th>
          {{ $yoy.DebtGrowth }}
        </th>
        <th>
          {{ $yoy.EquityGrowth }}
        </th>
        <th>
          {{ $yoy.BookValueGrowth }}
        </th>
        <th>
          {{ $yoy.DividendGrowth }}
        </th>
      </tr>
      {{end}}
      {{end}}
    </table>
  </body>
</html>
`
