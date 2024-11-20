import {EventsEmit, EventsOff, EventsOn, LogError} from '@wailsjs/runtime';
import {engin} from '@wailsjs/go/models';
import {Progress} from '@/components/ui/progress';
import {useEffect, useState} from 'react';
import {Button} from '@/components/ui/button';
import {GetConfigPath, SetConfigPath, Start} from '@wailsjs/go/gui/App';
import {Loader2} from 'lucide-react';
import {ScrollArea, ScrollBar} from '@/components/ui/scroll-area';
import {Accordion, AccordionContent, AccordionItem, AccordionTrigger} from '@/components/ui/accordion';
import Result = engin.ProgressResult;

// DebugInfoDisplay Displays debug info
function DebugInfoDisplay({debugInfo}: { debugInfo: Result }) {
    return (
        <div className="my-1 rounded-md px-2 py-0.5 text-white hover:bg-gray-800">
            <p className="inline">{debugInfo.time}&nbsp;</p>
            {debugInfo.msg && (!debugInfo.result && !debugInfo.err) &&
                <p className="inline">command={debugInfo.msg}&nbsp;</p>}
            {debugInfo.result && <p className="inline">output={debugInfo.result}&nbsp;</p>}
            {debugInfo.err && <p className="inline text-destructive">{debugInfo.err}&nbsp;</p>}
        </div>
    );
}

function HomePage() {
    const [progress, setProgress] = useState<Map<string, Result>>(new Map());
    const [debug, setDebug] = useState<Map<string, Array<Result>>>(new Map());
    const [config, setConfig] = useState<string>();
    const [run, setRun] = useState(false);

    // update Debug data info
    function debugSet(debugs: Map<string, Array<Result>>, result: Result): Map<string, Array<Result>> {
        const updated = new Map(debugs);
        const progressResults = updated.get(result.name) || [];
        updated.set(result.name, [...progressResults, result]);
        return updated;
    }

    useEffect(() => {
        // get config path
        GetConfigPath()
            .then(setConfig)
            .catch((err) => LogError(`Failed to get config path: ${err}`));

        // register for an event listener
        const handleProgressInfo = (r: Result) => {
            setProgress((prev) => {
                const updated = new Map(prev);
                updated.set(r.name, r);
                return updated;
            });
            setDebug((prev) => debugSet(prev, r));
        };

        EventsOn('progress_info', handleProgressInfo);

        // Clean up event listening
        return () => {
            EventsOff('progress_info');
        };
    }, []);

    const clickRun = () => {
        if (run) {
            EventsEmit('cancel_run');
        } else {
            setRun(true);
            setProgress(new Map());
            setDebug(new Map());

            Start('')
                .then(() => setRun(false))
                .catch((err) => {
                    LogError(`Failed to start process: ${err}`);
                    setRun(false);
                });
        }
    };

    const clickRetry = (clientIp: string) => {
        setRun(true);
        Start(clientIp)
            .then(() => setRun(false))
            .catch((err) => {
                LogError(`Failed to retry process: ${err}`)
                setRun(false);
            });
    };

    return (
        <div>
            <header className="sticky top-0 z-40 flex flex-row space-x-5 border-b bg-white px-8 py-8">
                <Button onClick={clickRun}>
                    {run ? (
                        <>
                            stop
                            <Loader2 className="mx-2 h-4 w-4 animate-spin"/>
                        </>
                    ) : (
                        'run'
                    )}
                </Button>

                <div className="flex flex-row space-x-2">
                    <Button onClick={() => SetConfigPath().then(setConfig)} disabled={run}>
                        select script
                    </Button>
                    <p className="py-2 text-slate-500">selected: {config}</p>
                </div>
            </header>

            <main className="container flex flex-col py-4">
                <Accordion type="single" collapsible>
                    <ul>
                        {Array.from(progress.values()).map((msg, index) => (
                            <li key={index} className="mb-4">
                                <div className="flex flex-row space-x-2 py-2">
                                    <p className="select-text text-slate-500">
                                        {index + 1}. {msg.name}
                                    </p>
                                    {run || msg.percent === 1 || (
                                        <p className="text-green-500 cursor-pointer"
                                           onClick={() => clickRetry(msg.name)}>
                                            retry
                                        </p>
                                    )}
                                </div>

                                <Progress value={msg.percent * 100}/>

                                <AccordionItem value={`index-${index}`}>
                                    <AccordionTrigger className="truncate">
                                        {msg.err ? <p className="inline text-destructive">{msg.err}</p> :
                                            <p>{msg.msg}</p>}
                                    </AccordionTrigger>
                                    <AccordionContent>
                    <pre className="mx-2 select-text rounded-md bg-slate-950">
                      <ScrollArea className="h-[400px] w-auto p-4">
                        <ul>
                          {(debug.get(msg.name) || []).map((debugInfo, idx) => (
                              <li key={idx}>
                                  <DebugInfoDisplay debugInfo={debugInfo}/>
                              </li>
                          ))}
                        </ul>
                        <ScrollBar orientation="horizontal"/>
                      </ScrollArea>
                    </pre>
                                    </AccordionContent>
                                </AccordionItem>
                            </li>
                        ))}
                    </ul>
                </Accordion>
            </main>
        </div>
    );
}

export default HomePage;
