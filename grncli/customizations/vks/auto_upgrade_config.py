from __future__ import annotations

from grncli.customizations.commands import BasicCommand, display_output
from grncli.customizations.vks.validators import validate_id


class SetAutoUpgradeConfigCommand(BasicCommand):
    NAME = 'set-auto-upgrade-config'
    DESCRIPTION = 'Configure auto-upgrade schedule for a cluster'
    ARG_TABLE = [
        {'name': 'cluster-id', 'help_text': 'Cluster ID', 'required': True},
        {'name': 'weekdays', 'help_text': 'Days of the week (e.g. Mon,Wed,Fri)', 'required': True},
        {'name': 'time', 'help_text': 'Time of day in 24h format HH:mm (e.g. 03:00)', 'required': True},
    ]

    def _run_main(self, parsed_args, parsed_globals):
        validate_id(parsed_args.cluster_id, 'cluster-id')
        client = self._session.create_client('vks')
        body = {
            'weekdays': parsed_args.weekdays,
            'time': parsed_args.time,
        }
        result = client.put(
            f'/v1/clusters/{parsed_args.cluster_id}/auto-upgrade-config',
            json=body,
        )
        display_output(result, parsed_globals)
        return 0


class DeleteAutoUpgradeConfigCommand(BasicCommand):
    NAME = 'delete-auto-upgrade-config'
    DESCRIPTION = 'Delete auto-upgrade config for a cluster'
    ARG_TABLE = [
        {'name': 'cluster-id', 'help_text': 'Cluster ID', 'required': True},
        {'name': 'force', 'help_text': 'Skip confirmation prompt',
         'action': 'store_true', 'default': False},
    ]

    def _run_main(self, parsed_args, parsed_globals):
        validate_id(parsed_args.cluster_id, 'cluster-id')
        client = self._session.create_client('vks')

        if not parsed_args.force:
            import sys
            response = input(
                "Are you sure you want to delete the auto-upgrade config? (yes/no): "
            )
            if response.strip().lower() != 'yes':
                sys.stdout.write("Delete cancelled.\n")
                return 0

        result = client.delete(
            f'/v1/clusters/{parsed_args.cluster_id}/auto-upgrade-config',
        )
        display_output(result, parsed_globals)
        return 0
