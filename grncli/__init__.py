# greenode-cli/grncli/__init__.py
from __future__ import annotations

import os

__version__ = '0.1.0'

# Data path for cli.json and other data files
_grncli_data_path = os.path.join(
    os.path.dirname(os.path.abspath(__file__)), 'data'
)
