# valuator
Stock Valuator

A package that gets input data about company filings and uses different valuation methods to come up with price range for stock prices.

Some useful measures:

Working Capital:
---------------
The difference between current assets and current liabilities
Current assets (CA):
    - Cash and cash equivalents
    - Accounts receivable
    - Short term investments
Current Liabilities (CL):
    - Accounts payable
    - Short term debt (due in the next year)
    - Interest and tax payable for the year

WC = CA - CL
Current Ratio = CA/CL

Cash Out:
--------

The amount of cash that can be taken out of a business over a period of time.
The cash out is the value being generated for the equity holder and is an
important component in the intrinsic value calculation

Cash out = Book Value increase (YoY) + Dividends paid out over that year

Discounted Cash Flow:
--------------------

The value of the cash generated for the equity holder up to a terminal date
discounted by a discount rate that decays the value of the cash being generated
over the period of time. The cash generated could also be derived from Free
Cash Flow (FCF) of the cash flow statement but it might hide some loss loss of
value to equity holders

The discount rate could be the risk-free treasury yield or a more company
specific rate by using the weighted average cost of capital (WACC)

DCF = cash out/(1-r)^year (series for the period of time being calculate)
