# Valuator
A Stock Valuator

A package that gets input data about company filings and uses different valuation methods to come up with price range for stock prices.

Design:

There are four main interfaces to the valuator:
  - Valuator
      The valuator is core interface of the package. The user gives a ticker to the interface. The valuator collects filing data upto the most recent filing (10K) and generates a bunch of measures, averages and Year over Year metrics. The interface then provides the user a number of valuation algorithms (ex: DCF) based on the collected metrics
  - Collector
      The collector interface collects filing data using a specific type of collector (ex: edgar) and populates the valuator with the filing data available. The collector could be used as a standalone interface to simply collect filing data and provide it to the user
  - Store
      The store interface is used by the valuator to store in-memory the filing data that has been collected and the measures and metrics that have been computed in the current iteration of the Valuator.
      The Store is an interface between the Valuator and any database that can store data that had already been evaluated. When the valuator is initialized, it reads the database that is in use by the valuator and populates the in-memory store with already evaluated data. The Valuator will then use the existing data and augment it with any more recent data available using the collector to update any metrics.
      The store and the collector could be used in conjuction directly for collecting filing data and storing it in a database without using the valuator.
  - Database
      Database is an interface for different underlying databases that could be used by the valuator.

Interfaces to various Metrics:
  - Measures
      Interface to the computed metrics from the data collected from filing. This is used as the basis for other calculated YoY metrics and averages
  - YoY
      This interface is a set of metrics generated from contiguous data available from the Measures interfaces. The interface reflects data YoY as calculated over the collected Measures
  - Averages
      This interface provides access to averages calculated for all the collected data and metrics. This forms the basis for projected valuations and intrinsic value calculations

REST interface:
  - Server
      The valuator provides a RESTful server that allows the user to start the valuator as a server and query information about a single ticker or a collection of tickers through a REST query.

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
Cash Flow (FCF) of the cash flow statement but it might hide some loss of
value to equity holders

The discount rate could be the risk-free treasury yield or a more company
specific rate by using the weighted average cost of capital (WACC)

An investor obtains this value and compares it to the current stock price and the return on the cost if it were to be invested in alternate investments.

DCF = cash out/(1-r)^year (series for the period of time being calculate)
