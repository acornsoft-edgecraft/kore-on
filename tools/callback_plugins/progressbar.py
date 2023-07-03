from ansible.plugins.callback import CallbackBase
from ansible import constants as C

class CallbackModule(CallbackBase):
    CALLBACK_VERSION = 2.0
    CALLBACK_TYPE = 'notification'
    CALLBACK_NAME = 'progressbar'
    CALLBACK_NEEDS_WHITELIST = True

    def __init__(self):
        super(CallbackModule, self).__init__()
        self.bar_length = 40

    def v2_playbook_on_play_start(self, play):
        self._display.banner("PLAY [{0}]".format(play.name))

    def v2_playbook_on_task_start(self, task, is_conditional):
        self._display.display("TASK [{0}]".format(task.get_name()), color=C.COLOR_SKIP)

    def v2_runner_on_ok(self, result, **kwargs):
        self._display.display("=> TASK OK", color=C.COLOR_OK)

    def v2_runner_on_failed(self, result, **kwargs):
        self._display.display("=> TASK FAILED", color=C.COLOR_ERROR)

    def v2_runner_on_skipped(self, result, **kwargs):
        self._display.display("=> TASK SKIPPED", color=C.COLOR_SKIP)
