"""
Portfolio Heat Calculator for Excel Python Integration

Purpose:
    Fast vectorized heat calculations using pandas.
    Validates trades against portfolio and bucket heat caps.

Usage in Excel:
    =PY("heat_calculator.check_heat_caps",
        xl("Positions[#All]"),
        TradeEntry!F5,
        TradeEntry!B8,
        Summary!B2,
        Summary!B5,
        Summary!B6)

Dependencies:
    - pandas
    - numpy

Author: Generated from newest-Interactive_TF_Workbook_Plan.md
"""

import pandas as pd
import numpy as np
from typing import Dict, List, Optional, Union


def portfolio_heat_after(positions_df: pd.DataFrame, add_r: float) -> float:
    """
    Calculates total portfolio heat (current open + proposed trade).

    Args:
        positions_df: DataFrame with columns [Ticker, Status, TotalOpenR, ...]
        add_r: Additional R dollars from proposed trade

    Returns:
        Total heat in dollars

    Example:
        >>> import pandas as pd
        >>> positions = pd.DataFrame({
        ...     'Ticker': ['AAPL', 'MSFT', 'GOOGL'],
        ...     'Status': ['Open', 'Open', 'Closed'],
        ...     'TotalOpenR': [75, 50, 100]
        ... })
        >>> portfolio_heat_after(positions, 75)
        200.0  # (75 + 50 from open positions) + 75 new = 200
    """
    if positions_df is None or positions_df.empty:
        return add_r

    # Filter to open positions only
    open_positions = positions_df[positions_df['Status'] != 'Closed']

    if open_positions.empty:
        return add_r

    # Sum TotalOpenR column
    current_heat = open_positions['TotalOpenR'].sum()

    # Handle NaN values
    if pd.isna(current_heat):
        current_heat = 0.0

    return float(current_heat) + float(add_r)


def bucket_heat_after(positions_df: pd.DataFrame, bucket: str, add_r: float) -> float:
    """
    Calculates bucket-specific heat (current open + proposed trade).

    Args:
        positions_df: DataFrame with columns [Bucket, Status, TotalOpenR, ...]
        bucket: Bucket name to filter by (e.g., "Tech/Comm")
        add_r: Additional R dollars from proposed trade

    Returns:
        Bucket heat in dollars

    Example:
        >>> import pandas as pd
        >>> positions = pd.DataFrame({
        ...     'Ticker': ['AAPL', 'MSFT', 'XOM'],
        ...     'Bucket': ['Tech/Comm', 'Tech/Comm', 'Energy/Materials'],
        ...     'Status': ['Open', 'Open', 'Open'],
        ...     'TotalOpenR': [75, 50, 60]
        ... })
        >>> bucket_heat_after(positions, 'Tech/Comm', 75)
        200.0  # (75 + 50) + 75 new = 200
    """
    if positions_df is None or positions_df.empty:
        return add_r

    # Filter to bucket + open positions
    bucket_positions = positions_df[
        (positions_df['Bucket'] == bucket) &
        (positions_df['Status'] != 'Closed')
    ]

    if bucket_positions.empty:
        return add_r

    # Sum TotalOpenR
    current_heat = bucket_positions['TotalOpenR'].sum()

    if pd.isna(current_heat):
        current_heat = 0.0

    return float(current_heat) + float(add_r)


def check_heat_caps(
    positions_df: pd.DataFrame,
    add_r: float,
    bucket: str,
    equity: float,
    port_cap_pct: float,
    bucket_cap_pct: float
) -> Dict[str, Union[bool, float]]:
    """
    Validates proposed trade against portfolio and bucket heat caps.

    Args:
        positions_df: DataFrame with position data
        add_r: Additional R dollars from proposed trade
        bucket: Bucket name for the proposed trade
        equity: Account equity (E)
        port_cap_pct: Portfolio heat cap as decimal (e.g., 0.04 for 4%)
        bucket_cap_pct: Bucket heat cap as decimal (e.g., 0.015 for 1.5%)

    Returns:
        Dictionary with keys:
            - portfolio_ok: bool (True if under cap)
            - bucket_ok: bool (True if under cap)
            - portfolio_heat: float (total heat after trade)
            - bucket_heat: float (bucket heat after trade)
            - portfolio_cap: float (cap in dollars)
            - bucket_cap: float (cap in dollars)
            - portfolio_pct: float (heat as % of equity)
            - bucket_pct: float (bucket heat as % of equity)

    Example:
        >>> result = check_heat_caps(positions_df, 75, 'Tech/Comm', 10000, 0.04, 0.015)
        >>> print(result)
        {
            'portfolio_ok': True,
            'bucket_ok': False,
            'portfolio_heat': 200.0,
            'bucket_heat': 200.0,
            'portfolio_cap': 400.0,
            'bucket_cap': 150.0,
            'portfolio_pct': 0.02,
            'bucket_pct': 0.02
        }
    """
    # Calculate heat values
    port_heat = portfolio_heat_after(positions_df, add_r)
    buck_heat = bucket_heat_after(positions_df, bucket, add_r)

    # Calculate caps
    port_cap = equity * port_cap_pct
    buck_cap = equity * bucket_cap_pct

    # Calculate percentages
    port_pct = port_heat / equity if equity > 0 else 0
    buck_pct = buck_heat / equity if equity > 0 else 0

    # Validate
    portfolio_ok = port_heat <= port_cap
    bucket_ok = buck_heat <= buck_cap

    return {
        'portfolio_ok': bool(portfolio_ok),
        'bucket_ok': bool(bucket_ok),
        'portfolio_heat': float(port_heat),
        'bucket_heat': float(buck_heat),
        'portfolio_cap': float(port_cap),
        'bucket_cap': float(buck_cap),
        'portfolio_pct': float(port_pct),
        'bucket_pct': float(buck_pct)
    }


def get_open_positions_summary(positions_df: pd.DataFrame) -> Dict[str, any]:
    """
    Returns summary statistics of open positions.

    Args:
        positions_df: DataFrame with position data

    Returns:
        Dictionary with summary stats

    Example:
        >>> summary = get_open_positions_summary(positions_df)
        >>> print(summary)
        {
            'total_positions': 3,
            'total_heat': 185.0,
            'buckets': {'Tech/Comm': 125.0, 'Energy/Materials': 60.0},
            'avg_position_size': 61.67,
            'largest_position': 75.0
        }
    """
    if positions_df is None or positions_df.empty:
        return {
            'total_positions': 0,
            'total_heat': 0.0,
            'buckets': {},
            'avg_position_size': 0.0,
            'largest_position': 0.0
        }

    # Filter to open positions
    open_pos = positions_df[positions_df['Status'] != 'Closed'].copy()

    if open_pos.empty:
        return {
            'total_positions': 0,
            'total_heat': 0.0,
            'buckets': {},
            'avg_position_size': 0.0,
            'largest_position': 0.0
        }

    # Calculate summary stats
    total_positions = len(open_pos)
    total_heat = float(open_pos['TotalOpenR'].sum())

    # Bucket breakdown
    bucket_heat = open_pos.groupby('Bucket')['TotalOpenR'].sum().to_dict()

    # Position stats
    avg_size = float(open_pos['TotalOpenR'].mean())
    largest = float(open_pos['TotalOpenR'].max())

    return {
        'total_positions': int(total_positions),
        'total_heat': total_heat,
        'buckets': bucket_heat,
        'avg_position_size': avg_size,
        'largest_position': largest
    }


def calculate_max_position_size(
    equity: float,
    risk_pct: float,
    current_heat: float,
    port_cap_pct: float,
    bucket_heat: float,
    bucket_cap_pct: float
) -> Dict[str, float]:
    """
    Calculates maximum allowable position size given current heat.

    Args:
        equity: Account equity
        risk_pct: Risk % per unit (e.g., 0.0075 for 0.75%)
        current_heat: Current portfolio heat
        port_cap_pct: Portfolio cap %
        bucket_heat: Current bucket heat
        bucket_cap_pct: Bucket cap %

    Returns:
        Dictionary with max position sizes:
            - max_r_portfolio: Max R allowed by portfolio cap
            - max_r_bucket: Max R allowed by bucket cap
            - max_r_combined: Min of the two (actual max)
            - max_units: Max units at current risk %

    Example:
        >>> max_size = calculate_max_position_size(
        ...     equity=10000,
        ...     risk_pct=0.0075,
        ...     current_heat=300,
        ...     port_cap_pct=0.04,
        ...     bucket_heat=100,
        ...     bucket_cap_pct=0.015
        ... )
        >>> print(max_size)
        {
            'max_r_portfolio': 100.0,  # 400 cap - 300 current
            'max_r_bucket': 50.0,      # 150 cap - 100 current
            'max_r_combined': 50.0,    # Limited by bucket
            'max_units': 0.67          # 50 / 75 = 0.67 units
        }
    """
    # Calculate caps
    port_cap = equity * port_cap_pct
    buck_cap = equity * bucket_cap_pct
    r_per_unit = equity * risk_pct

    # Calculate remaining room
    max_r_portfolio = max(0, port_cap - current_heat)
    max_r_bucket = max(0, buck_cap - bucket_heat)

    # Combined limit is the stricter one
    max_r_combined = min(max_r_portfolio, max_r_bucket)

    # Convert to units
    max_units = max_r_combined / r_per_unit if r_per_unit > 0 else 0

    return {
        'max_r_portfolio': float(max_r_portfolio),
        'max_r_bucket': float(max_r_bucket),
        'max_r_combined': float(max_r_combined),
        'max_units': float(max_units)
    }


def validate_position_sizing(
    positions_df: pd.DataFrame,
    ticker: str,
    bucket: str,
    proposed_r: float,
    equity: float,
    risk_pct: float,
    port_cap_pct: float,
    bucket_cap_pct: float
) -> Dict[str, any]:
    """
    Comprehensive validation of a proposed position.

    Combines heat checks, max size calculation, and warnings.

    Args:
        positions_df: Current positions DataFrame
        ticker: Proposed ticker symbol
        bucket: Bucket for proposed trade
        proposed_r: R dollars for proposed trade
        equity: Account equity
        risk_pct: Risk % per unit
        port_cap_pct: Portfolio cap %
        bucket_cap_pct: Bucket cap %

    Returns:
        Dictionary with validation results and recommendations

    Example:
        >>> result = validate_position_sizing(
        ...     positions_df, 'AAPL', 'Tech/Comm', 75, 10000, 0.0075, 0.04, 0.015
        ... )
        >>> if result['valid']:
        ...     print("Trade approved")
        ... else:
        ...     print(f"Blocked: {result['reasons']}")
    """
    # Run heat check
    heat_check = check_heat_caps(
        positions_df, proposed_r, bucket, equity, port_cap_pct, bucket_cap_pct
    )

    # Calculate current heat
    current_port_heat = portfolio_heat_after(positions_df, 0)
    current_bucket_heat = bucket_heat_after(positions_df, bucket, 0)

    # Calculate max allowable
    max_size = calculate_max_position_size(
        equity, risk_pct, current_port_heat, port_cap_pct,
        current_bucket_heat, bucket_cap_pct
    )

    # Determine validity
    valid = heat_check['portfolio_ok'] and heat_check['bucket_ok']

    # Build reasons list
    reasons = []
    if not heat_check['portfolio_ok']:
        reasons.append(f"Portfolio heat would be {heat_check['portfolio_pct']:.1%} "
                      f"(cap: {port_cap_pct:.1%})")
    if not heat_check['bucket_ok']:
        reasons.append(f"Bucket heat would be {heat_check['bucket_pct']:.1%} "
                      f"(cap: {bucket_cap_pct:.1%})")

    return {
        'valid': valid,
        'reasons': reasons,
        'heat_check': heat_check,
        'max_size': max_size,
        'current_portfolio_heat': current_port_heat,
        'current_bucket_heat': current_bucket_heat,
        'ticker': ticker,
        'bucket': bucket,
        'proposed_r': proposed_r
    }


# Test function for development
def _test_heat_calculator():
    """
    Test function - not called by Excel.
    Run this in Python to verify calculator works.
    """
    # Create test positions
    positions = pd.DataFrame({
        'Ticker': ['AAPL', 'MSFT', 'XOM', 'GOOGL'],
        'Bucket': ['Tech/Comm', 'Tech/Comm', 'Energy/Materials', 'Tech/Comm'],
        'Status': ['Open', 'Open', 'Open', 'Closed'],
        'TotalOpenR': [75, 50, 60, 100],
        'UnitsOpen': [1, 1, 1, 0]
    })

    print("Testing Heat Calculator...")
    print("\nCurrent Positions:")
    print(positions[positions['Status'] == 'Open'][['Ticker', 'Bucket', 'TotalOpenR']])

    # Test portfolio heat
    port_heat = portfolio_heat_after(positions, 0)
    print(f"\nCurrent Portfolio Heat: ${port_heat:.2f}")

    # Test bucket heat
    tech_heat = bucket_heat_after(positions, 'Tech/Comm', 0)
    print(f"Tech/Comm Bucket Heat: ${tech_heat:.2f}")

    # Test proposed trade
    print("\n--- Proposed Trade ---")
    print("Ticker: NVDA, Bucket: Tech/Comm, R: $75")

    result = check_heat_caps(
        positions,
        add_r=75,
        bucket='Tech/Comm',
        equity=10000,
        port_cap_pct=0.04,
        bucket_cap_pct=0.015
    )

    print(f"\nPortfolio Heat: ${result['portfolio_heat']:.2f} / ${result['portfolio_cap']:.2f} "
          f"({result['portfolio_pct']:.1%}) - {'✅ OK' if result['portfolio_ok'] else '❌ OVER'}")
    print(f"Bucket Heat: ${result['bucket_heat']:.2f} / ${result['bucket_cap']:.2f} "
          f"({result['bucket_pct']:.1%}) - {'✅ OK' if result['bucket_ok'] else '❌ OVER'}")

    # Test max size calculation
    max_size = calculate_max_position_size(
        equity=10000,
        risk_pct=0.0075,
        current_heat=port_heat,
        port_cap_pct=0.04,
        bucket_heat=tech_heat,
        bucket_cap_pct=0.015
    )

    print(f"\nMax Allowable Position:")
    print(f"  By Portfolio Cap: ${max_size['max_r_portfolio']:.2f}")
    print(f"  By Bucket Cap: ${max_size['max_r_bucket']:.2f}")
    print(f"  Actual Max: ${max_size['max_r_combined']:.2f} ({max_size['max_units']:.2f} units)")


if __name__ == "__main__":
    # Run test if executed directly
    _test_heat_calculator()
