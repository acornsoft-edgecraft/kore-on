# Ansible

- version: v2.11

## Ansible 필수 패키지 / Plugin

- community.docker

```sh
ansible-galaxy collection install community.docker
```

- ansible_mitogen plugin
  Mitogen for Ansible은 Ansible을 위해 완전히 재설계된 UNIX 연결 계층 및 모듈 런타임입니다.
  기대 효과: ansible 실행 속도 개선

```sh
## ansible.cfg 샘플
## Mitogen for Ansible is a completely redesigned UNIX connection layer and module runtime for Ansible.
strategy_plugins = ./tools/mitogen-0.3.2/ansible_mitogen/plugins/strategy
strategy = mitogen_linear

## 소스 다운로드
$ wget https://github.com/mitogen-hq/mitogen/archive/refs/tags/v0.3.2.tar.gz

## 참고 사이트
참고: https://github.com/mitogen-hq/mitogen/releases/tag/v0.3.2
```