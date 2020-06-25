# Cetus (주의: 아래 문서는 컨셉으로 아직 구현체가 없습니다)

쿠버네티스는 확장성이 높지만, 많은 세부 기술들(service mesh, image registry, package manager, serverless, monitoring, logging, cni plugin, container runtime 등)이 별도 프로젝트로 산재해 있어 러닝커브가 크고 관리가 쉽지 않습니다. 때문에 [2019 KubeCon & CloudNativeCon 설문 결과](https://thenewstack.io/ux-is-kubernetes-biggest-short-term-challenge/?fbclid=IwAR1Olut6i5Ekf4TQ3-QQ7P5jEaYNuan3s73ndzV8HOXf6Yc06Hu_QjtIkxk) 많은 운영자들과 개발자들이 유저 경험(UX)이 해결해야할 중요 단기 과제라고 답했습니다.

Cetus는 쿠버네티스를 빌딩 블록 삼아 개발자와 운영자들에게 더 나은 UX를 제공하는 통합 컨테이너 환경(Integrated Container Environment)으로써 운영 측면에서는 쿠버네티스에서 기본적으로 제공하는 기능 이외에 멀티테넌시(CaaS multitenancy model 2), 인증/인가, OIDC, 모니터링, 로깅, 트레이싱, 다양한 배포(A/B 테스팅, 카나리, Blue/Green), 서킷 브레이커, 이미지 레지스트리, 승인 제어(Admission Control), 정책 제어(Policy Control), 쉬운 노드/스토리지 추가/삭제, 백업/복구를 제공하며, 개발 측면에서는 로컬 개발과 배포 환경의 간극을 줄이고, 다른 개발자가 작성한 서비스를 의존성으로 추가하여 쉽게 사용할 수 있으며, 직관적인 배포 설정(cetus.yaml), 로깅 및 트레이싱 및 로컬 UI를 통해 버그를 쉽게 추적하여 수정할 수 있도록 합니다.

## 암묵적인 컨텍스트

- 아래 정보는 프로젝트의 .cetus/config에 저장된다.
- cluster namespace (user or group)
- image registry
- .cetus/cache 디렉토리 아래에 build에 관련된 캐시가 저장된다.

## 명시적인 설정

- 아래 정보는 \<project root\>/cetus.yaml 파일에 개발자가 작성한다.
- images: build 명령어 실행시 빌드될 이미지 정보
- services: 배포될 서비스들(k8s deployment + service로 구성)
- dependencies: 의존성을 가지는 외부 서비스(git, path 참조를 지원)
- profiles: patch(rfc6902), replace를 지원하여 cetus.yaml의 값들을 수정하여 배포 할 수 있다.
- dev: 개발시 필요한 설정들
- test: 배포 후 확인 커맨드 실행. 실패할 경우 기존 버전으로 rollback 된다.
- vars: file, env, stdin 입력을 받아 cetus.yaml 내에서 사용 할 수 있다.
- jobs: 
- ingress: 외부에서 domain으로 서비스에 접근할 수 있다.

## 초기화

cetus는 배포 환경(인프라)와 개발 환경(프로젝트)를 각각 초기화 할 수 있다. 로컬에서 배포하려면 cetus init 만으로 개발을 시작 할 수 있지만, 배포하려면 인프라가 필요하다. 인프라를 초기화하면 로컬 쿠버네티스 클러스터에 cetus 운영에 필요한 커스텀 리소스(CRD) 및 컨트롤 플레인을 설치하며, 제공된 대시보드 UI에서 github 혹은 gitlab 처럼 그룹/유저별로 네임스페이스를 만들어 로컬 프로젝트에 네임스페이스를 추가 할 수 있다.

### 인프라 초기화

- cetus install (options)
    - 설치시 대화형으로 어드민 계정/비밀번호를 입력한다.
    - --domain 옵션에 도메인을 제공하면 제공된 URL로 cetus 대시보드에 접근 가능하다.
    - --domain 을 지정하지 않더라도 아래와 같이 임의의 도메인이 제공된다.
        - 예) https://5c7110be.cetus.dev
    - 대시보드에서는 사용자, 인증, 네임스페이스, 모니터링,

### 프로젝트 초기화

- cetus init
- cetus namespace add staging https://5c7110be.cetus.dev/wickedev/bookinfo
    - 배포시 클러스터 내에서는 wickedev--bookinfo k8s 네임스페이스를 가짐
- cetus namespace add prod https://5c7110be.cetus.dev/demo/bookinfo
    - 배포시 클러스터 내에서는 demo--bookinfo k8s 네임스페이스를 가짐

## 개발

- cetus dev (alias)
- 해당 어플리케이션을 제외하고, 의존성이 있는 서비스(DB, MQ, 마이크로서비스 등)을 로컬 쿠버네티스 클러스에 배포한다.
- (alias) 인자 없이 실행 할 경우 로컬에 배포되며, 클러스터를 지정하여 배포할 수 있다.
- 배포된 각 서비스를 로컬 프로세스 0번 포트에 바인딩 하여 임의로 할당한 뒤, 현재 프로세스 환경 변수에 바인딩 한다.
    - 가령 postgresql 서비스가 로컬 30423 포트에 바인딩 되었다면, POSTGRESQL_SERVICE_HOST는 localhost, POSTGRESQL_SERVICE_PORT는 30423이다.
    - 해당 환경 변수는 어플리케이션을 실행하는 프로세스에만 바인딩 된다.

## 배포

- cetus deploy [alias]
- cetus deploy staging
- cetus deploy prod

- deploy 호출 시 인증 정보를 요구 할 수 있다. 인증은 user + password 혹은 token(base65 encoded x509)
- deploy 호출 전에 이미지 publish(build, push), 호출 후에는 test가 불려지며 실패할 경우 rollback이 수행된다.
- publish, build, push는 git 과 .cetus/cache 을 참조하여 캐시, 빌드, 버전을 결정한다.
    - 버전은 (branch)-(8 length git hash) 이다. 예시 master-1we09e41
    - 기본 동작은 git에 커밋된 코드만 빌드하지만 --force 플래그를 사용하면 커밋되지 않은 코드를 포함하며 버전 +(revision)이 붙는다. 예시 master-1we09e41+1
- cetus test
    - cetus.yaml에 명시한 테스트를 수행한다.

## 롤백

- cetus status [alias]
    - k8s resources 및 그 상태 (uptime, age, restart, status)
- cetus history [alias]
    - 지금까지 배포된 버전 목록
- cetus rollback [alias] (version)
    - 버전을 지정하지 않을 경우 대화식으로 선택 제공

## 디버깅

- cetus debug [alias] (pod) (-c [container])
    - 최근 컨테이너가 죽었을 당시 pod log 및 events 출력
- cetus exec [alias] (pod) (-c [container]) -- [exec]
    - 예) cetus exec staging -- ls -al
    - 컨테이너를 지정하지 않을 경우 대화식으로 선택 제공
- cetus attach [alias] (pod) (-c [container])
    - 예) cetus attach staging
    - 컨테이너를 지정하지 않을 경우 대화식으로 선택 제공
- cetus logs (options) [alias] (pod) (-c [container])
    - 예) cetus logs -f staging
    - 컨테이너를 지정하지 않을 경우 namespace의 전체 로그 (주입된 사이드카 제외)
- cetus curl (-n [alias]) [url]
    - 클러스터 URL(svc.cluster.local) 혹은 도메인으로 클러스터 내부에 curl을 수행
    - 예) cetus curl http://bookinfo.wickedev--bookinfo.svc.cluster.local
    - 예) cetus curl https://staging.bookinfo.5c7110be.cetus.dev
    - 예) cetus curl -n wickedev/bookinfo BOOKINFO_SERVICE_HOST
    - 예) cetus curl -n staging BOOKINFO_SERVICE_HOST:8080

## 기타

- cetus ui
    - dev로 배포된 서비스 상태 확인 또는 관리 할 수 있는 UI를 브라우저에 띄운다.
